package aabf

func onNetPackGR(c *C, flag, msgid uint32, msgdata []byte) {
	if c.dir != nil {
		c.SendData(flag, msgid, msgdata)
		return
	}

	mi := minfoById.Load()[MID(msgid)]
	if mi != nil {
		onNetPack(c, flag, msgid, msgdata)
		return
	}

	dest := LB.Get(MID(msgid))
	if dest != nil {
		dest.SendData(flag, msgid, msgdata)
		return
	}

	c.Disconnect()
	Log.W("onNetPackGR => Invalid client message " + c.CInfo)
}
