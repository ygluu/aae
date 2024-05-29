package aabf

import (
	"sync"
)

func exeMsg(senderQ *q, exefunc, callback func(), key string) {
	defer Recover()

	exefunc()

	if callback != nil {
		senderQ.push(nil, 1, "", 0, nil, callback, nil)
	}
}

type asynG struct {
	G
}

func (a *asynG) OnMsg(node *qnode) {
	exeMsg(node.senderQ, node.proc1, node.proc2, node.sp)

	asyn.Lock()
	defer asyn.Unlock()

	if a.getQ().Count() == 0 {
		delete(asyn.asynGByKey, node.sp)
	}
}

type asynMgr struct {
	G
	sync.Mutex
	asynGByKey map[string]*asynG
	index      int
	asynGS     []*asynG
}

func (a *asynMgr) onMsg(node *qnode) {
	key := node.sp
	if key == "" {
		go exeMsg(node.senderQ, node.proc1, node.proc2, node.sp)
		return
	}

	a.Lock()
	defer a.Unlock()

	c := a.asynGByKey[key]
	if c == nil {
		c = a.asynGS[a.index]
		a.index++
		if a.index >= len(a.asynGS) {
			a.index = 0
		}
		for _, fc := range a.asynGS {
			if c == fc {
				continue
			}
			if fc.getQ().Count() < c.getQ().Count() {
				c = fc
			}
		}
	}

	c.q.push(nil, 0, key, 0, node.senderQ, node.proc1, node.proc2)
}

var asyn *asynMgr

func setAsynGoCount(v int) {
	var newGos []*asynG
	for i := 1; i < v; i++ {
		c := &asynG{}
		RegModule(c)
		newGos = append(newGos, c)
	}
	newGos = append(newGos, asyn.asynGS...)
	asyn.asynGS = newGos
}

func init() {
	asyn = &asynMgr{
		asynGByKey: make(map[string]*asynG),
	}
	RegModule(asyn)
	setAsynGoCount(20)
}

func SetAsynGoCountByKey(v int) {
	setAsynGoCount(v - 20)
}
