package session

import (
	"sync"
)

var sessionDictPool = sync.Pool{
	New: func() interface{} {
		return new(Dict)
	},
}

// AcquireDict acquire new Dict
func AcquireDict() *Dict {
	return sessionDictPool.Get().(*Dict)
}

// ReleaseDict release Dict
func ReleaseDict(d *Dict) {
	d.Reset()
	sessionDictPool.Put(d)
}
