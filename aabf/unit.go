package aabf

import (
	"fmt"
	"os"
	"runtime"
	"sync/atomic"
	"unsafe"
)

// msg id
type MID uint32

// connect id
type CID uint64

var isInitTime = true

func IsInitTime() bool {
	return isInitTime
}

func FileExists(name string) bool {
	info, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func DirExists(name string) bool {
	info, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func Panic(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	Log.Exp(msg)
	panic(msg)
}

func Recover() {
	err := recover()
	if err == nil {
		return
	}

	stackTrace := make([]byte, 1024*1024*2)
	length := runtime.Stack(stackTrace, false)

	Log.Exp(fmt.Sprintf("%v", err) + string(stackTrace[:length]))
}

type Map[K comparable, V any] struct {
	v map[K]V
	p unsafe.Pointer
}

func (m *Map[K, V]) New() map[K]V {
	return make(map[K]V)
}

func (m *Map[K, V]) Clone() map[K]V {
	ret := make(map[K]V)

	for k, v := range m.v {
		ret[k] = v
	}

	return ret
}

func (m *Map[K, V]) Store(v map[K]V) {
	m.v = v
	atomic.StorePointer(&m.p, unsafe.Pointer(&m.v))
}

func (m *Map[K, V]) Load() map[K]V {
	if m.p == nil {
		m.Store(m.New())
	}
	return *(*map[K]V)(atomic.LoadPointer(&m.p))
}
