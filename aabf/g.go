package aabf

// G: goroutine obj

import (
	"unsafe"
)

const flag_stop = uint64(0xFFFFFFFFFFFFFFFF)

type gi interface {
	Init()
	PostMsg(msg MsgI)
	AsynOper(exefunc, callback func(), key string)
	getQ() *q
	setQ(q *q)
	getType() int
	Name() string
	setName(name string)
	Ready()
	Close()
	Track(v int)
	Timer() *timer
	setTimer(t *timer)
	setG(g gi)
	GID() uint64
	setGID(v uint64)

	stop()
	isEnd() bool
	isRuning() bool
	getTrack() int
	onStart()
	onStop()
	onMsg(node *qnode)
	getMode() int
}

// Coroutine
type G struct {
	M
	isexe    bool
	isend    bool
	isruning bool
	track    int
	mode     int
	goId     uint64
}

func (g *G) Init() {
	g.M.Init()
	g.q = newQ()
	g.isexe = true
	g.t = newTimer(g.q)
}

func (g *G) Ready() {
}

func (g *G) Close() {
	g.t.Clear()
}

func (g *G) Timer() *timer {
	return g.t
}

func (g *G) getMode() int {
	return g.mode
}

func (g *G) Track(v int) {
	g.track = v
}

func (g *G) getTrack() int {
	return g.track
}

func (g *G) getType() int {
	return 2
}

func (g *G) isEnd() bool {
	return g.isend
}

func (g *G) isRuning() bool {
	return g.isruning
}

func (g *G) onStart() {
	g.isruning = true
}

func (g *G) onStop() {
	g.isend = true
}

func (g *G) GID() uint64 {
	return g.goId
}

func (g *G) setGID(v uint64) {
	g.goId = v
}

type peerCIDToObj func(mid MID, cid CID) (unsafe.Pointer, int)

var PeerCIDToObj peerCIDToObj = func(mid MID, cid CID) (unsafe.Pointer, int) {
	return nil, 1
}

func (g *G) onMsg(node *qnode) {
	defer Recover()
	switch node.flag {
	case 3:
		sminfo := (*sminfo)(unsafe.Pointer(node.senderQ))
		msg := node.msg
		h := msg.GetH()

		if (h == nil) || (sminfo.paramCnt != 2) {
			sminfo.method((*eface)(unsafe.Pointer(&msg)).data, nil)
			break
		}

		obj, ret := PeerCIDToObj(msg.GetId(), h.GetPeerCID())
		if ret == 0 {
			sminfo.method((*eface)(unsafe.Pointer(&msg)).data, obj)
		}
		break
	case 2:
		g.t.onTimer(node.np, node.proc1, node.proc2)
		break
	case 1:
		node.proc1()
		break
	}
}

func (g *G) stop() {
	g.q.push(nil, flag_stop, "", 0, nil, nil, nil)
}

type mainG struct {
	G
}

func (g *mainG) Init() {
	g.G.Init()
	g.mode = 1
}

func (g *mainG) onStart() {
	g.G.onStart()
	doRegModule()
}

var MainG = &mainG{}

func init() {
	MainG.Init()
}
