package sqlite3

import (
	"sync"
	"time"

	"github.com/fasthttp/session"
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

// NewStore new sqlite3 store
func NewStore(sessionID []byte, data *session.Dict) *Store {
	store := acquireStore()
	store.Init(sessionID, data)

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
