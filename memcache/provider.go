package memcache

import (
	"sync"

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
		config:         new(Config),
		memCacheClient: new(memcache.Client),
	}
}

// Init init provider config
func (mcp *Provider) Init(lifeTime int64, cfg session.ProviderConfig) error {
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

	mcp.memCacheClient = memcache.New(mcp.config.ServerList...)
	mcp.memCacheClient.MaxIdleConns = mcp.config.MaxIdleConns
	mcp.maxLifeTime = lifeTime

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

// ReadStore read session store by session id
func (mcp *Provider) ReadStore(sessionID []byte) (session.Storer, error) {
	store := NewStore(sessionID)

	item := acquireItem()
	item, err := mcp.memCacheClient.Get(mcp.getMemCacheSessionKey(sessionID))
	if err == nil { // Exist
		err := mcp.config.UnSerializeFunc(item.Value, store.GetDataPointer())
		if err != nil {
			return nil, err
		}

	} else if err == memcache.ErrCacheMiss { // Not exist
		err = nil
	}

	releaseItem(item)

	return store, err
}

// Regenerate regenerate session
func (mcp *Provider) Regenerate(oldID, newID []byte) (session.Storer, error) {
	store := NewStore(newID)

	oldKey := mcp.getMemCacheSessionKey(oldID)
	newKey := mcp.getMemCacheSessionKey(newID)

	oldItem := acquireItem()
	oldItem, err := mcp.memCacheClient.Get(oldKey)
	if err == nil { // Exist
		newItem := acquireItem()
		newItem.Key = newKey
		newItem.Value = oldItem.Value
		newItem.Expiration = oldItem.Expiration

		if err = mcp.memCacheClient.Set(newItem); err != nil {
			return nil, err
		}

		if err = mcp.memCacheClient.Delete(oldKey); err != nil {
			return nil, err
		}

		err := mcp.config.UnSerializeFunc(newItem.Value, store.GetDataPointer())
		if err != nil {
			return nil, err
		}

		releaseItem(newItem)

	} else if err == memcache.ErrCacheMiss { // Not exist
		err = nil
	}

	releaseItem(oldItem)

	return store, err
}

// Destroy destroy session by sessionID
func (mcp *Provider) Destroy(sessionID []byte) error {
	return mcp.memCacheClient.Delete(mcp.getMemCacheSessionKey(sessionID))
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
	session.Register(ProviderName, provider)
}
