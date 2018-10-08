package mysql

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

// NewStore new mysql store
func NewStore(sessionID []byte) *Store {
	store := acquireStore()
	store.Init(sessionID)

	return store
}

// Save save store
func (ms *Store) Save(ctx *fasthttp.RequestCtx) error {
	defer releaseStore(ms)

	data := ms.GetData()
	value, err := provider.config.SerializeFunc(*data)
	if err != nil {
		return err
	}

	_, err = provider.db.updateBySessionID(ms.GetSessionID(), value, time.Now().Unix())

	return err
}
