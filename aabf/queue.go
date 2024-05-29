package aabf

import (
	"sync"
	"sync/atomic"
)

// queue obj

type qnode struct {
	next    *qnode
	flag    uint64
	senderQ *q
	sp      string
	np      uint64
	msg     MsgI
	hs      []H
	proc1   func()
	proc2   func()
}

type q struct {
	sync.Mutex
	count     int
	head      *qnode
	last      *qnode
	frees     *qnode
	freeCount int
	waitCount int32
	c         chan bool
}

func newQ() *q {
	ret := &q{
		c: make(chan bool, 30000),
	}
	return ret
}

func (q *q) Count() int {
	return 0
}

func (q *q) push(msg MsgI, flag uint64, sp string, np uint64, senderQ *q, proc1, proc2 func()) {
	q.Lock()
	defer q.Unlock()

	var n *qnode
	if q.frees == nil {
		n = &qnode{}
	} else {
		n = q.frees
		q.frees = n.next
		q.freeCount--
	}

	n.msg = msg
	n.flag = flag
	n.sp = sp
	n.np = np
	n.senderQ = senderQ
	n.proc1 = proc1
	n.proc2 = proc2

	n.next = nil
	q.count++
	if q.last == nil {
		q.head = n
		q.last = n
	} else {
		q.last.next = n
		q.last = n
	}

	if atomic.LoadInt32(&q.waitCount) > 0 {
		q.c <- true
	}
}

func (q *q) pushMsg(msg MsgI, p uint64) {
	q.Lock()
	defer q.Unlock()

	var n *qnode
	if q.frees == nil {
		n = &qnode{}
	} else {
		n = q.frees
		q.frees = n.next
		q.freeCount++
	}

	n.msg = msg
	n.flag = 100
	n.np = p

	n.next = nil
	q.count++
	if q.last == nil {
		q.head = n
		q.last = n
	} else {
		q.last.next = n
		q.last = n
	}

	if atomic.LoadInt32(&q.waitCount) > 0 {
		q.c <- true
	}
}

func (q *q) pushMsgTos(msg MsgI, hs []H) {
	q.Lock()
	defer q.Unlock()

	var n *qnode
	if q.frees == nil {
		n = &qnode{}
	} else {
		n = q.frees
		q.frees = n.next
		q.freeCount++
	}

	n.msg = msg
	n.flag = 101
	n.hs = hs

	n.next = nil
	q.count++
	if q.last == nil {
		q.head = n
		q.last = n
	} else {
		q.last.next = n
		q.last = n
	}

	if atomic.LoadInt32(&q.waitCount) > 0 {
		q.c <- true
	}
}

func (q *q) pushBrd(msg MsgI, p uint64) {
	q.Lock()
	defer q.Unlock()

	var n *qnode
	if q.frees == nil {
		n = &qnode{}
	} else {
		n = q.frees
		q.frees = n.next
		q.freeCount++
	}

	n.msg = msg
	n.flag = 102
	n.np = p

	n.next = nil
	q.count++
	if q.last == nil {
		q.head = n
		q.last = n
	} else {
		q.last.next = n
		q.last = n
	}

	if atomic.LoadInt32(&q.waitCount) > 0 {
		q.c <- true
	}
}

func (q *q) pushs(head, last *qnode, count int) {
	q.Lock()
	defer q.Unlock()

	last.next = nil
	q.count += count
	if q.last == nil {
		q.head = head
		q.last = last
	} else {
		q.last.next = head
		q.last = last
	}

	if atomic.LoadInt32(&q.waitCount) > 0 {
		q.c <- true
	}
}

func (q *q) freeNodes(head, last *qnode, count int) {
	q.Lock()
	defer q.Unlock()
	last.next = q.frees
	q.frees = head
	q.freeCount += count
}

func (q *q) getFrees() (ret *qnode) {
	q.Lock()
	defer q.Unlock()

	if q.freeCount == 0 {
		return &qnode{}
	}

	if q.freeCount < 1000 {
		ret = q.frees
		q.frees = nil
		return
	}

	ret = q.frees
	for i := 1; i < 999; i++ {
		q.frees = q.frees.next
	}
	last := q.frees
	q.frees = q.frees.next
	last.next = nil

	return
}

func (q *q) doPops() (head, last *qnode, count int) {
	q.Lock()
	defer q.Unlock()

	count = q.count
	head = q.head
	last = q.last

	q.count = 0
	q.head = nil
	q.last = nil

	return
}

func (q *q) doPop(node *qnode) bool {
	q.Lock()
	defer q.Unlock()

	head := q.head
	if head == nil {
		return false
	}

	q.count--
	q.head = q.head.next
	if q.head == nil {
		q.last = nil
	}

	*node = *head
	node.next = nil

	head.next = q.frees
	q.frees = head

	return true
}

func (q *q) pop(node *qnode) (ret bool) {
	for {
		ret = q.doPop(node)
		if !ret {
			atomic.AddInt32(&q.waitCount, 1)
			<-q.c
			atomic.AddInt32(&q.waitCount, -1)
			continue
		}
		return
	}

	return
}

func (q *q) pops() (head, last *qnode, count int) {
	for {
		head, last, count = q.doPops()
		if head == nil {
			atomic.AddInt32(&q.waitCount, 1)
			<-q.c
			atomic.AddInt32(&q.waitCount, -1)
			continue
		}
		return
	}

	return
}

func (q *q) newSubQ() *subq {
	return &subq{owner: q}
}

type subq struct {
	owner *q
	count int
	head  *qnode
	last  *qnode
	frees *qnode
}

func (q *subq) getFree() (ret *qnode) {
	if q.frees == nil {
		q.frees = q.owner.getFrees()
	}
	ret = q.frees
	q.frees = q.frees.next
	return
}

func (q *subq) push(msg MsgI, flag uint64, sp string, np uint64, senderQ *q, proc1, proc2 func()) {
	n := q.getFree()

	n.msg = msg
	n.flag = flag
	n.sp = sp
	n.np = np
	n.senderQ = senderQ
	n.proc1 = proc1
	n.proc2 = proc2

	if q.last == nil {
		q.head = n
		q.last = n
	} else {
		q.last.next = n
		q.last = n
	}

	q.count++
}

func (q *subq) submit() {
	q.owner.pushs(q.head, q.last, q.count)
	q.head = nil
	q.last = nil
	q.count = 0
}
