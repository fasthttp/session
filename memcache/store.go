package memcache

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/fasthttp/session"
	"github.com/valyala/fasthttp"
)

// session memCache store

// NewMemCacheStore new default memCache store
func NewMemCacheStore(sessionID string) *Store {
	memCacheStore := &Store{}
	memCacheStore.Init(sessionID, make(map[string]interface{}))
	return memCacheStore
}

// NewMemCacheStoreData new memCache store data
func NewMemCacheStoreData(sessionID string, data map[string]interface{}) *Store {
	memCacheStore := &Store{}
	memCacheStore.Init(sessionID, data)
	return memCacheStore
}

// Store store struct
type Store struct {
	session.Store
}

// Save save store
func (mcs *Store) Save(ctx *fasthttp.RequestCtx) error {

	value, err := provider.config.SerializeFunc(mcs.GetAll())
	if err != nil {
		return err
	}

	return provider.memCacheClient.Set(&memcache.Item{
		Key:        provider.getMemCacheSessionKey(mcs.GetSessionID()),
		Value:      value,
		Expiration: int32(provider.maxLifeTime),
	})
}
