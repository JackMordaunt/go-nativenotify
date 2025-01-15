package nativenotify

import (
	"sync"
	"sync/atomic"
)

var (
	nextID    atomic.Int64
	callbacks sync.Map
)

func callbacksTake(m *sync.Map, id string) (Callback, bool) {
	v, ok := m.LoadAndDelete(id)
	if !ok {
		return nil, false
	}
	cb, ok := v.(Callback)
	if !ok {
		return nil, false
	}
	return cb, true
}

func callbacksPut(m *sync.Map, id string, fn Callback) {
	m.Store(id, fn)
}
