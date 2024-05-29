package aabf

type cinfo struct {
	lb    *lb
	cname string
}

type ainfo struct {
	ci    *cinfo
	sname string
	saddr string
}

var cinfoByName = make(map[string]*cinfo)
var ainfoByAddr = make(map[string]*ainfo)

func addSvcIns(cname, sname, saddr string, startid uint32, alloAutoConnect bool) {
	ci := cinfoByName[cname]
	if ci == nil {
		ci = &cinfo{
			lb:    newLb(),
			cname: cname,
		}
		cinfoByName[cname] = ci
	}

	ai := ainfoByAddr[saddr]
	if ai == nil {
		ai = &ainfo{
			ci:    ci,
			sname: sname,
			saddr: saddr,
		}
		ainfoByAddr[saddr] = ai
	}

	if NetSets.IsGate || (alloAutoConnect && (startId > startid)) {
		Connect(saddr)
	}
}

func delSvcIns(cname, sname, saddr string) {
	delete(ainfoByAddr, saddr)
}
