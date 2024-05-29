package aabf

import (
	"container/list"
	"math"
	"sync/atomic"
	"time"
)

const min_interval = 50
const sleep_time = 20

var timerIdCount uint64 = 0

type timer struct {
	ownerQ *q
	ids    map[uint64]bool
}

func (t *timer) Add(life, interval uint32, onTimer, onDie func()) {
	if (life > 0) && (life < min_interval) {
		life = min_interval
	}
	if interval < min_interval {
		interval = min_interval
	}
	id := atomic.AddUint64(&timerIdCount, 1)
	t.ids[id] = true
	tc.q.push(nil, (uint64(life)<<32)+uint64(interval), "", id, t.ownerQ, onTimer, onDie)
}

func (t *timer) Delete(id uint64) {
	delete(t.ids, id)
}

func (t *timer) Clear() {
	t.ids = make(map[uint64]bool)
}

func (t *timer) onTimer(id uint64, onTimer, onDie func()) {
	if !t.ids[id] {
		return
	}
	if onTimer != nil {
		onTimer()
	}
	if onDie != nil {
		onDie()
	}
}

func newTimer(ownerQ *q) *timer {
	return &timer{
		ownerQ: ownerQ,
		ids:    make(map[uint64]bool),
	}
}

type tinfo struct {
	id        uint64
	life      uint32
	interval  uint32
	starttime int64
	lasttime  int64
	count     uint32
	ownerQ    *q
	onTimer   func()
	onDie     func()
}

type timerMgr struct {
	G
	list       *list.List
	isTcRuning bool
	isTcEnd    bool
}

func (a *timerMgr) doTime() {
	for a.isTcRuning {
		time.Sleep(time.Millisecond * sleep_time)
		a.q.push(nil, 200, "", 0, nil, nil, nil)
	}
	a.isTcEnd = true
	Log.I("Timer is end")
}

func (a *timerMgr) Init() {
	a.G.Init()
	a.list = list.New()
	a.isTcRuning = true
	go a.doTime()
}

func (a *timerMgr) Close() {
	a.isTcRuning = false
	Log.I("Timer is closing...")
	for {
		if a.isTcEnd {
			break
		}
		time.Sleep(time.Millisecond * min_interval)
	}
}

func (a *timerMgr) check() {
	currtime := time.Now().UnixMilli()
	mp := make(map[*q]*subq)

	for oe := a.list.Front(); oe != nil; oe = oe.Next() {
		ti := oe.Value.(*tinfo)

		long := uint32(currtime - ti.lasttime)
		if long < ti.interval {
			continue
		}

		onTimer := ti.onTimer
		onDie := ti.onDie
		isDie := false
		if ti.life > 0 {
			if ti.count == 0 {
				onTimer = nil
			} else {
				ti.count--
			}
			isDie = (uint32(currtime-ti.starttime) >= ti.life)
		}

		if onTimer != nil {
			ti.lasttime = currtime
		}

		if !isDie {
			onDie = nil
		}

		subq := mp[ti.ownerQ]
		if subq == nil {
			subq = ti.ownerQ.newSubQ()
			mp[ti.ownerQ] = subq
		}
		subq.push(nil, 2, "", ti.id, nil, onTimer, onDie)

		if isDie {
			a.list.Remove(oe)
		}
	}

	for _, q := range mp {
		q.submit()
	}
}

func (a *timerMgr) onMsg(node *qnode) {
	defer Recover()

	if node.flag == 200 {
		a.check()
		return
	}

	life := uint32(node.flag >> 32)
	interval := uint32(node.flag & 0xFFFFFFFF)

	nti := &tinfo{
		id:       node.np,
		life:     life,
		interval: interval,
		count:    uint32(math.Round(float64(life / interval))),
		lasttime: time.Now().UnixMilli(),
		ownerQ:   node.senderQ,
		onTimer:  node.proc1,
		onDie:    node.proc2,
	}

	if nti.count == 0 {
		nti.count = 1
	}
	// if uint64(node.flag)%node.np2 > 0 {
	// 	nti.count++
	// }
	nti.starttime = nti.lasttime

	ne := a.list.PushBack(nti)

	for oe := a.list.Front(); oe != nil; oe = oe.Next() {
		if (oe.Next() == ne) || (ne == oe) {
			return
		}
		oti := oe.Value.(*tinfo)
		if nti.interval > oti.interval {
			a.list.InsertAfter(oe, ne)
			return
		}
	}
}

var tc = &timerMgr{}

func init() {
	RegModule(tc)
}
