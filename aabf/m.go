package aabf

// M: module obj

type mi interface {
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
}

type M struct {
	g    gi
	name string
	q    *q
	t    *timer
}

func (m *M) Init() {

}

func (m *M) setG(g gi) {
	m.g = g
}

func (m *M) Track(v int) {
	m.g.Track(v)
}

func (m *M) Timer() *timer {
	return m.t
}

func (m *M) setTimer(t *timer) {
	m.t = t
}

func (m *M) Ready() {

}

func (m *M) Close() {

}

func (m *M) Name() string {
	return m.name
}

func (m *M) setName(name string) {
	m.name = name
}

func (m *M) getType() int {
	return 1
}

func (m *M) getQ() *q {
	return m.q
}

func (m *M) setQ(q *q) {
	m.q = q
}

func (m *M) GID() uint64 {
	return m.g.GID()
}

func (m *M) Broadcast(msg MsgI) {
	sndr.q.pushBrd(msg, 0)
}

func (m *M) PostMsg(msg MsgI) {
	sndr.q.pushMsg(msg, 0)
}

func (m *M) PostTo(msg MsgI, h H) {
	sndr.q.pushMsgTos(msg, []H{h})
}

func (m *M) PostTos(msg MsgI, hs []H) {
	sndr.q.pushMsgTos(msg, hs)
}

func (m *M) BroadcastToAll(msg MsgI) {
	sndr.q.pushBrd(msg, 1)
}

func (m *M) PostMsgToAll(msg MsgI) {
	sndr.q.pushMsg(msg, 1)
}

func (m *M) AsynOper(exefunc, callback func(), key string) {
	asyn.q.push(nil, 0, key, 0, m.getQ(), exefunc, callback)
}
