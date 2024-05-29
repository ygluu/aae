package aabf

var startId uint32

func getStartId() uint32 {
	return startId
}

type GlobalID uint64

func NewGlobalID() GlobalID {
	return GlobalID(uint64(startId) << 32)
}

func (g GlobalID) Inc() uint64 {
	if g == 0 {
		g = GlobalID(uint64(getStartId()) << 32)
	}
	ret := g
	g++
	return uint64(ret)
}
