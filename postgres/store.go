package postgres

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

// NewStore new postgres store
func NewStore(sessionID []byte) *Store {
	store := acquireStore()
	store.Init(sessionID)

	return store
}

// Save save store
func (ps *Store) Save(ctx *fasthttp.RequestCtx) error {
	defer releaseStore(ps)

	data := ps.GetAll()
	value, err := provider.config.SerializeFunc(data)
	if err != nil {
		return err
	}

	_, err = provider.db.updateBySessionID(ps.GetSessionID(), value, time.Now().Unix())

	return err
}
