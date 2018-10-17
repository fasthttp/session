package memory

import (
	"sync"
	"time"

	"github.com/fasthttp/session"
)

var provider = NewProvider()

// NewProvider new memory provider
func NewProvider() *Provider {
	return &Provider{
		config:      new(Config),
		memoryDB:    new(session.Dict),
		maxLifeTime: 0,

		storePool: sync.Pool{
			New: func() interface{} {
				return new(Store)
			},
		},
	}
}

func (mp *Provider) acquireStore(sessionID []byte) *Store {
	store := mp.storePool.Get().(*Store)
	store.Init(sessionID)

	return store
}

func (mp *Provider) releaseStore(store *Store) {
	store.Reset()
	mp.storePool.Put(store)
}

// Init init provider configuration
func (mp *Provider) Init(lifeTime int64, cfg session.ProviderConfig) error {
	if cfg.Name() != ProviderName {
		return errInvalidProviderConfig
	}

	mp.config = cfg.(*Config)
	mp.maxLifeTime = lifeTime

	return nil
}

// Get get session store by id
func (mp *Provider) Get(sessionID []byte) (session.Storer, error) {
	currentStore := mp.memoryDB.GetBytes(sessionID)
	if currentStore != nil {
		return currentStore.(*Store), nil
	}

	newStore := mp.acquireStore(sessionID)
	mp.memoryDB.SetBytes(sessionID, newStore)

	return newStore, nil
}

// Put put store into the pool.
//
// In Memory provider, only put again the store into the pool when destroy the session
func (mp *Provider) Put(store session.Storer) {}

// Regenerate regenerate session
func (mp *Provider) Regenerate(oldID, newID []byte) (session.Storer, error) {
	var store *Store

	val := mp.memoryDB.GetBytes(oldID)
	if val != nil {
		store = val.(*Store)
		store.SetSessionID(newID)
		mp.memoryDB.SetBytes(newID, store)
		mp.memoryDB.DelBytes(oldID)
	} else {
		store = mp.acquireStore(newID)
		mp.memoryDB.SetBytes(newID, store)
	}

	return store, nil
}

// Destroy destroy session by sessionID
func (mp *Provider) Destroy(sessionID []byte) error {
	val := mp.memoryDB.GetBytes(sessionID)
	if val != nil {
		mp.releaseStore(val.(*Store))
	}

	mp.memoryDB.DelBytes(sessionID)

	return nil
}

// Count session values count
func (mp *Provider) Count() int {
	return len(mp.memoryDB.D)
}

// NeedGC need gc
func (mp *Provider) NeedGC() bool {
	return true
}

// GC session garbage collection
func (mp *Provider) GC() {
	for _, kv := range mp.memoryDB.D {
		if time.Now().Unix() >= (kv.Value.(*Store).lastActiveTime + mp.maxLifeTime) {
			mp.Destroy(kv.Key)
		}
	}
}

// register session provider
func init() {
	session.Register(ProviderName, provider)
}
