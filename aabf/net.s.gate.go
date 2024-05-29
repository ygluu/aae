package aabf

func svc2ClientPack(c *C, flag, msgid uint32, msgdata []byte) {
	const errfmt = "AAGE => Send invalid message packets(%d), flag:%d, msgid:%d, len:%d"

	datalen := len(msgdata)
	if datalen <= 2 {
		Log.E(errfmt, 1, flag, msgid, len(msgdata))
		return
	}

	count := int(uint32(msgdata[1])<<8) + int(msgdata[0])
	cidlen := 2 + count*8
	if cidlen >= datalen {
		Log.E(errfmt, 2, flag, msgid, len(msgdata))
		return
	}

	di := 2
	cids := make([]uint64, count)
	for i := 0; i < count; i++ {
		cids[i] = uint64(msgdata[di]) +
			uint64(msgdata[di+1])<<8 +
			uint64(msgdata[di+2])<<16 +
			uint64(msgdata[di+3])<<24 +
			uint64(msgdata[di+4])<<32 +
			uint64(msgdata[di+5])<<40 +
			uint64(msgdata[di+6])<<48 +
			uint64(msgdata[di+7])<<56
		di += 8
	}

	dests := make([]*C, count)
	netLock.RLock()
	for i := 0; i < count; i++ {
		dests[i] = connOfClt[CID(cids[i])]
	}
	netLock.RUnlock()

	senddata := msgdata[cidlen:]
	for _, dest := range dests {
		if dest != nil {
			dest.SendData(flag, msgid, senddata)
		}
	}
}

func brd2ClientPack(c *C, flag, msgid uint32, msgdata []byte) {
	for _, dest := range connOfClt {
		dest.SendData(flag, msgid, msgdata)
	}

	i := 0
	netLock.RLock()
	dests := make([]*C, len(connOfClt))
	for _, dest := range connOfClt {
		dests[i] = dest
	}
	netLock.RUnlock()

	for _, dest := range dests {
		if dest != nil {
			dest.SendData(flag, msgid, msgdata)
		}
	}
}

func onNetPackGS(c *C, flag, msgid uint32, msgdata []byte) {
	if flag&0x4 > 0 {
		brd2ClientPack(c, flag, msgid, msgdata)
		return
	}

	if flag&0x8 > 0 {
		brd2ClientPack(c, flag, msgid, msgdata)
		return
	}

	onNetPack(c, flag, msgid, msgdata)
}
