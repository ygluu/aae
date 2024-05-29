package aabf

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

var ms []mi
var gs []gi
var mgs []mi
var IsRuning = false

func execG1(g gi) {
	node := &qnode{}
	q := g.getQ()
	for {
		if q.pop(node) {
			if node.flag == flag_stop {
				return
			}
			g.onMsg(node)
		}
	}
}

func execG2(g gi) {
	var head *qnode
	var last *qnode
	count := 0
	q := g.getQ()
	for {
		head, last, count = q.pops()
		if head == nil {
			continue
		}
		node := head
		for node != nil {
			if node.flag == flag_stop {
				return
			}
			g.onMsg(node)
			node = node.next
		}
		q.freeNodes(head, last, count)
	}
}

func CurrGoID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func execG(g gi) {
	g.setGID(CurrGoID())

	Log.I("The goroutine is starting >> %d - %s", g.GID(), g.Name())

	g.onStart()
	if g.getMode() == 0 {
		execG1(g)
	} else {
		execG2(g)
	}
	g.onStop()

	Log.I("The goroutine is end >> %d - %s", g.GID(), g.Name())
}

func Exec() {
	Log.I("The ggo system is starting...")

	for _, mg := range mgs {
		mg.Init()
	}

	for _, m := range ms {
		m.setG(MainG)
		m.setQ(MainG.getQ())
		m.setTimer(MainG.Timer())
	}

	for _, mg := range mgs {
		mg.Ready()
	}

	for _, g := range gs {
		go execG(g)
	}

	time.Sleep(time.Millisecond * 100)

	for _, g := range gs {
		if g.isRuning() {
			continue
		}
		Log.I("Waiting runing for goroutine >> %d - %s", g.GID(), g.Name())
		for !g.isRuning() {
			time.Sleep(time.Millisecond * 10)
		}
	}

	IsRuning = true
	Log.I("The ggo system is ready")

	execG(MainG)
}

func Stop() {
	Log.I("The ggo system is stoping...")

	for _, mg := range mgs {
		mg.Close()
	}

	time.Sleep(time.Millisecond * 200)
	for _, g := range gs {
		if g.getQ().Count() == 0 {
			continue
		}
		Log.I("Waiting empty for queue >> %d - %s", g.GID(), g.Name())
		for g.getQ().Count() > 0 {
			time.Sleep(time.Millisecond * 10)
		}
	}

	for _, g := range gs {
		g.stop()
	}

	time.Sleep(time.Millisecond * 200)
	for _, g := range gs {
		if g.isEnd() {
			continue
		}
		Log.I("Waiting end for goroutine >> %d - %s", g.GID(), g.Name())
		for !g.isEnd() {
			time.Sleep(time.Millisecond * 10)
		}
	}

	Log.I("The ggo system is stoped")

	time.Sleep(time.Second * 3)
}

func RegModule(m mi) {
	t := reflect.TypeOf(m)
	m.setName(t.Elem().String() + fmt.Sprintf("(%p)", m))

	mt := m.getType()
	switch mt {
	case 1:
		ms = append(ms, m)
		mgs = append(mgs, m)
		break
	case 2:
		v := reflect.ValueOf(m)
		g, ok := v.Interface().(gi)
		if !ok {
			panic("The goroutine module is invalid")
		}
		gs = append(gs, g)
		mgs = append(mgs, g)
		break
	default:
		panic("The module is invalid")
	}
}

func doRegModule() {
	for _, m := range ms {
		regMethods(m)
	}

	for _, m := range gs {
		regMethods(m)
	}
}
