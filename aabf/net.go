package aabf

import (
	"net"
	"sync"
)

type onpack = func(c *C, flag, msgid uint32, msgdata []byte)

func goRecv(c *C) {

}

// C: connect
type C struct {
	packCount   uint64
	id          CID
	c           net.Conn
	peerCluster string
	peerSvc     string
	peerAddr    string
	onpack      onpack
	IsSvr       bool
	lb          *lb

	dir   *C
	CInfo string
}

func (c *C) SendData(flag, msgid uint32, msgdata []byte) {
}

func (c *C) Send(msg MsgI, peer CID) {

}

func (c *C) Sends(msg MsgI, peers []CID) {

}

func (c *C) Disconnect() {
}

type netEvent func(c *C)

type netSets struct {
	CName string
	SName string
	SAddr string

	AllowAutoConnect bool

	IsGate bool

	MaxSByte int
	MaxRByte int

	OnCltConnect    netEvent
	OnCltDisconnect netEvent
	OnSvrConnect    netEvent
	OnSvrDisconnect netEvent
}

func onCltConnect(c *C) {
}

func onCltDisconnect(c *C) {

}

func onSvrConnect(c *C) {
}

func onSvrDisconnect(c *C) {

}

func init() {
	NetSets.AllowAutoConnect = true

	NetSets.OnCltConnect = onCltConnect
	NetSets.OnCltDisconnect = onCltDisconnect
	NetSets.OnSvrConnect = onSvrConnect
	NetSets.OnSvrDisconnect = onSvrDisconnect
}

var netLock sync.RWMutex
var NetSets = &netSets{}
var listener net.Listener
var connOfSvc = make(map[CID]*C)
var connOfClt = make(map[CID]*C)
var lclose bool
var netGID GlobalID

type msg_net_connect struct {
	Msg
	c *C
}

type msg_net_disconnect struct {
	Msg
	c *C
}

type netMgr struct {
	G
}

func (n *netMgr) Ready() {
	n.Timer().Add(
		0, 1000,
		func() {

		},
		nil,
	)
}

func (n *netMgr) OnConnect(msg *msg_net_connect) {
	c := msg.c
	c.id = CID(netGID.Inc())
	if c.IsSvr {
		connOfSvc[c.id] = c
	} else {
		connOfClt[c.id] = c
	}

	go goRecv(c)
}

func (n *netMgr) OnDisconnect(msg *msg_net_disconnect) {
	c := msg.c
	if c.IsSvr {
		delete(connOfSvc, c.id)
	} else {
		delete(connOfSvc, c.id)
	}
}

func (n *netMgr) PostConnect(c *C) {
	msg := &msg_net_connect{c: c}
	n.PostMsg(msg)
}

func (n *netMgr) PostDisconnect(c *C) {
	msg := &msg_net_disconnect{c: c}
	n.PostMsg(msg)
}

var netmgr = &netMgr{}

func init() {
	RegModule(netmgr)
}

func accept(l net.Listener) {
	for {
		c, err := l.Accept()
		if (c == nil) || (err != nil) {
			lclose = true
			return
		}

		onpack := onNetPack
		if NetSets.IsGate {
			onpack = onNetPackGR
		}
		conn := &C{
			c:      c,
			onpack: onpack,
			IsSvr:  true,
			lb:     LB,
		}
		netmgr.PostConnect(conn)
	}
}

func Listen(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	listener = l
	go accept(l)
	return nil
}

func Connect(addr string) (H, error) {
	lb := LB
	ai := ainfoByAddr[addr]
	if ai != nil {
		lb = ai.ci.lb
	}

	c, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	onpack := onNetPack
	if NetSets.IsGate {
		onpack = onNetPackGS
	}
	conn := &C{
		c:      c,
		lb:     lb,
		onpack: onpack,
	}

	netmgr.PostConnect(conn)
	return &h{localC: conn}, nil
}
