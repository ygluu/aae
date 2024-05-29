package aabf

// loadbalan

type ring struct {
	index uint32
	cs    []*C
}

type tmidsByC map[*C][]MID
type tringByMID map[MID]*ring
type tcsByMID map[MID][]*C

type lb struct {
	midsByC   tmidsByC
	ringByMID tringByMID
	csByMID   tcsByMID
}

type lbs []*lb

var LBS lbs

var LB = newLb()

func newLb() *lb {
	return &lb{
		midsByC:   make(tmidsByC),
		ringByMID: make(tringByMID),
		csByMID:   make(tcsByMID),
	}
}

func init() {
	LBS = lbs{LB}
}

func addC(cs []*C, v *C) {
	for _, c := range cs {
		if c == v {
			return
		}
	}
	cs = append(cs, v)
}

func (l *lb) Add(c *C, mids []MID) {
	l.midsByC[c] = mids

	for _, mid := range mids {
		addC(l.csByMID[mid], c)
		r := l.ringByMID[mid]
		if r != nil {
			addC(r.cs, c)
		}
	}
}

func delC(cs []*C, v *C) {
	for i, dest := range cs {
		if dest != v {
			continue
		}

		ln := len(cs)
		if ln == 1 {
			cs = nil
		} else if i == 0 {
			cs = cs[1:]
		} else if i == len(cs)-1 {
			cs = cs[:i]
		} else {
			cs = append(cs[:i], cs[i+1:]...)
		}

		return
	}
}

func (l *lb) Del(c *C) {
	mids := l.midsByC[c]
	delete(l.midsByC, c)

	for _, mid := range mids {
		delC(l.csByMID[mid], c)
		r := l.ringByMID[mid]
		if r != nil {
			delC(r.cs, c)
		}
	}
}

func (l *lb) Get(mid MID) *C {
	r := l.ringByMID[mid]
	if r == nil {
		return nil
	}

	ln := len(r.cs)
	if ln == 0 {
		return nil
	}

	r.index++
	return r.cs[int(r.index%uint32(ln))]
}

func (l *lb) Gets(mid MID) []*C {
	return l.csByMID[mid]
}
