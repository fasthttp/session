package memory

import (
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

var storePool = sync.Pool{
	New: func() interface{} {
		return new(Store)
	},
}

func acquireStore() *Store {
	return storePool.Get().(*Store)
}

func releaseStore(store *Store) {
	store.Reset()
	storePool.Put(store)
}

// NewStore new memory store
func NewStore(sessionID []byte) *Store {
	memStore := acquireStore()
	memStore.Init(sessionID)

	return memStore
}

// Save save store
func (ms *Store) Save(ctx *fasthttp.RequestCtx) error {
	ms.lock.Lock()
	ms.lastActiveTime = time.Now().Unix()
	ms.lock.Unlock()

	return nil
}
