package aabf

import (
	"reflect"
	"unsafe"
)

type sender struct {
	G
}

type msg_connect struct {
	Msg
	c    *C
	lb   *lb
	mids []MID
}

type msg_disconnect struct {
	Msg
	c  *C
	lb *lb
}

func init() {
	regMsg(&msg_connect{}, false, false)
	regMsg(&msg_disconnect{}, false, false)
}

func (a *sender) OnCconnect(msg *msg_connect) {
	msg.lb.Add(msg.c, msg.mids)
}

func (a *sender) OnDisconnect(msg *msg_disconnect) {
	msg.lb.Del(msg.c)
}

func (a *sender) onMsg(node *qnode) {
	if node.flag < 100 {
		a.g.onMsg(node)
		return
	}

	defer Recover()

	msg := node.msg
	if msg == nil {
		Log.W("Send empty message >> %v", node)
		return
	}

	mt := reflect.TypeOf(msg)
	minfo := minfoByType.Load()[mt.Name()]
	if minfo == nil {
		Log.W("An undefined message has occurred >> %v -- %v", mt.Name(), msg)
		return
	}

	for _, sminfo := range minfo.sminfos {
		sminfo.recvQ.push(msg, 3, "", 0, (*q)(unsafe.Pointer(sminfo)), nil, nil)
	}

	if !minfo.proto {
		return
	}

	if node.flag == 0 {
		if node.np == 1 {
			for _, lb := range LBS {
				c := lb.Get(minfo.id)
				if c != nil {
					c.Send(msg, 0)
				}
			}
			return
		}
		c := LB.Get(minfo.id)
		if c != nil {
			c.Send(msg, 0)
		}
		return
	}

	if node.flag == 2 {
		if node.np == 1 {
			for _, lb := range LBS {
				cs := lb.Gets(minfo.id)
				for _, c := range cs {
					c.Send(msg, 0)
				}
			}
			return
		}
		cs := LB.Gets(minfo.id)
		for _, c := range cs {
			c.Send(msg, 0)
		}
	}

	ln := len(node.hs)

	if ln == 1 {
		c := LB.Get(minfo.id)
		if c != nil {
			c.Send(msg, node.hs[0].GetPeerCID())
		}
		return
	}

	peerCIDsByC := make(map[*C][]CID)
	for _, h := range node.hs {
		c := h.GetLocalC()
		pcids := peerCIDsByC[c]
		if pcids == nil {
			pcids = []CID{}
			peerCIDsByC[c] = pcids
		}
		pcids = append(pcids, h.GetPeerCID())
	}
	for c, pcids := range peerCIDsByC {
		if c == nil {
			return
		}
		if pcids == nil {
			return
		}
		c.Sends(msg, pcids)
	}
}

var sndr = &sender{}

func init() {
	RegModule(sndr)
}
