package memcache

import (
	"sync"

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

// NewStore new memCache store
func NewStore(sessionID []byte, data *session.Dict) *Store {
	store := acquireStore()
	store.Init(sessionID, data)

	return store
}

// Save save store
func (mcs *Store) Save(ctx *fasthttp.RequestCtx) error {
	defer releaseStore(mcs)

	data := mcs.GetAll()
	value, err := provider.config.SerializeFunc(&data)

	if err != nil {
		return err
	}

	item := acquireItem()
	item.Key = provider.getMemCacheSessionKey(mcs.GetSessionID())
	item.Value = value
	item.Expiration = int32(provider.maxLifeTime)

	err = provider.memCacheClient.Set(item)

	releaseItem(item)

	return err
}
