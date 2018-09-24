package memory

import (
	"time"

	"github.com/valyala/fasthttp"
)

// NewStore new memory store
func NewStore(sessionID []byte) *Store {
	memStore := new(Store)
	memStore.Init(sessionID, nil)

	return memStore
}

// Save save store
func (ms *Store) Save(ctx *fasthttp.RequestCtx) error {
	ms.lock.Lock()
	ms.lastActiveTime = time.Now().Unix()
	ms.lock.Unlock()

	return nil
}
