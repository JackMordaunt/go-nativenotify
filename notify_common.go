package nativenotify

import (
	"encoding/hex"
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

// take returns the first element.
func take[S ~[]E, E any](s *S) (e E, ok bool) {
	if len(*s) == 0 {
		return e, false
	}
	defer func() { *s = (*s)[1:] }()
	return (*s)[0], true
}

// decode from hex.
func decode(s string) string {
	b, _ := hex.DecodeString(s)
	return string(b)
}

// encode to hex.
func encode(s string) string {
	return hex.EncodeToString([]byte(s))
}
