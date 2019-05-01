package memcache

import (
	"math"
	"sync"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/fasthttp/session"
	"github.com/valyala/bytebufferpool"
)

var (
	provider = NewProvider()
	encrypt  = session.NewEncrypt()
	itemPool = sync.Pool{
		New: func() interface{} {
			return new(memcache.Item)
		},
	}
)

func acquireItem() *memcache.Item {
	return itemPool.Get().(*memcache.Item)
}

func releaseItem(item *memcache.Item) {
	if item != nil {
		item.Key = ""
		item.Value = nil
		item.Expiration = 0

		itemPool.Put(item)
	}
}

// NewProvider new memcache provider
func NewProvider() *Provider {
	return &Provider{
		config: new(Config),
		db:     new(memcache.Client),

		storePool: sync.Pool{
			New: func() interface{} {
				return new(Store)
			},
		},
	}
}

func (mcp *Provider) acquireStore(sessionID []byte, expiration time.Duration) *Store {
	store := mcp.storePool.Get().(*Store)
	store.Init(sessionID, expiration)

	return store
}

func (mcp *Provider) releaseStore(store *Store) {
	store.Reset()
	mcp.storePool.Put(store)
}

// Init init provider config
func (mcp *Provider) Init(expiration time.Duration, cfg session.ProviderConfig) error {
	if cfg.Name() != ProviderName {
		return errInvalidProviderConfig
	}

	mcp.config = cfg.(*Config)

	// config check
	if len(mcp.config.ServerList) == 0 {
		return errConfigServerListEmpty
	}
	if mcp.config.MaxIdleConns <= 0 {
		return errConfigMaxIdleConnsZero
	}

	// init config serialize func
	if mcp.config.SerializeFunc == nil {
		mcp.config.SerializeFunc = encrypt.MSGPEncode
	}
	if mcp.config.UnSerializeFunc == nil {
		mcp.config.UnSerializeFunc = encrypt.MSGPDecode
	}

	mcp.db = memcache.New(mcp.config.ServerList...)
	mcp.db.MaxIdleConns = mcp.config.MaxIdleConns
	if expiration/time.Second > math.MaxInt32 {
		return errExpirationIsTooBig
	}
	mcp.expiration = expiration

	return nil
}

// get memcache session key, prefix:sessionID
func (mcp *Provider) getMemCacheSessionKey(sessionID []byte) string {
	key := bytebufferpool.Get()
	key.SetString(mcp.config.KeyPrefix)
	key.WriteString(":")
	key.Write(sessionID)

	keyStr := key.String()

	bytebufferpool.Put(key)

	return keyStr
}

// Get read session store by session id
func (mcp *Provider) Get(sessionID []byte) (session.Storer, error) {
	var store *Store
	key := mcp.getMemCacheSessionKey(sessionID)

	item, err := mcp.db.Get(key)
	if err != nil && err != memcache.ErrCacheMiss {
		return nil, err
	}

	if item != nil { // Exist
		store = mcp.acquireStore(sessionID, time.Duration(item.Expiration)*time.Second)

		err := mcp.config.UnSerializeFunc(store.DataPointer(), item.Value)
		if err != nil {
			return nil, err
		}
	} else {
		store = mcp.acquireStore(sessionID, mcp.expiration)
	}

	releaseItem(item)

	return store, nil
}

// Put put store into the pool.
func (mcp *Provider) Put(store session.Storer) {
	mcp.releaseStore(store.(*Store))
}

// Regenerate regenerate session
func (mcp *Provider) Regenerate(oldID, newID []byte) (session.Storer, error) {
	store := mcp.acquireStore(newID, mcp.expiration)

	oldKey := mcp.getMemCacheSessionKey(oldID)
	newKey := mcp.getMemCacheSessionKey(newID)

	oldItem, err := mcp.db.Get(oldKey)
	if err != nil && err != memcache.ErrCacheMiss {
		return nil, err
	}

	if oldItem != nil { // Exist
		newItem := acquireItem()
		newItem.Key = newKey
		newItem.Value = oldItem.Value
		// Expiration can be converted safely from time.Duration to in32 because
		// mcp.expiration has been checked to be safely castable in Init method.
		newItem.Expiration = int32(mcp.expiration / time.Second)

		if err = mcp.db.Set(newItem); err != nil {
			return nil, err
		}

		if err = mcp.db.Delete(oldKey); err != nil {
			return nil, err
		}

		err := mcp.config.UnSerializeFunc(store.DataPointer(), newItem.Value)
		if err != nil {
			return nil, err
		}

		releaseItem(newItem)
	}

	releaseItem(oldItem)

	return store, nil
}

// Destroy destroy session by sessionID
func (mcp *Provider) Destroy(sessionID []byte) error {
	key := mcp.getMemCacheSessionKey(sessionID)
	return mcp.db.Delete(key)
}

// Count session values count
func (mcp *Provider) Count() int {
	return 0
}

// NeedGC not need gc
func (mcp *Provider) NeedGC() bool {
	return false
}

// GC session memcache provider not need garbage collection
func (mcp *Provider) GC() {}

// register session provider
func init() {
	err := session.Register(ProviderName, provider)
	if err != nil {
		panic(err)
	}
}
