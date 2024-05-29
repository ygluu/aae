package aabf

// connect handle

type H interface {
	SetLocalC(c *C)
	GetLocalC() *C
	SetPeerCID(cid CID)
	GetPeerCID() CID
}

type h struct {
	localC  *C
	peerCID CID
}

func NewH(c *C, cid CID) H {
	return &h{localC: c, peerCID: cid}
}

func (h *h) SetLocalC(c *C) {
	h.localC = c
}

func (h *h) GetLocalC() *C {
	return h.localC
}

func (h *h) SetPeerCID(cid CID) {
	h.peerCID = cid
}

func (h *h) GetPeerCID() CID {
	return h.peerCID
}
