package memcache

import (
	"math"
	"sync"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/fasthttp/session"
	"github.com/valyala/bytebufferpool"
)

var itemPool = &sync.Pool{
	New: func() interface{} {
		return new(memcache.Item)
	},
}

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

// New new memcache provider
func New(cfg Config) (*Provider, error) {
	// config check
	if len(cfg.ServerList) == 0 {
		return nil, errConfigServerListEmpty
	}
	if cfg.MaxIdleConns <= 0 {
		return nil, errConfigMaxIdleConnsZero
	}

	// init config serialize func
	if cfg.SerializeFunc == nil {
		cfg.SerializeFunc = session.MSGPEncode
	}
	if cfg.UnSerializeFunc == nil {
		cfg.UnSerializeFunc = session.MSGPDecode
	}

	db := memcache.New(cfg.ServerList...)
	db.MaxIdleConns = cfg.MaxIdleConns

	p := &Provider{
		config: cfg,
		db:     db,
	}

	return p, nil
}

func (p *Provider) getMemCacheSessionKey(sessionID []byte) string {
	key := bytebufferpool.Get()
	key.SetString(p.config.KeyPrefix)
	key.WriteString(":")
	key.Write(sessionID)

	keyStr := key.String()

	bytebufferpool.Put(key)

	return keyStr
}

// Get read session store by session id
func (p *Provider) Get(store *session.Store) error {
	key := p.getMemCacheSessionKey(store.GetSessionID())

	item, err := p.db.Get(key)
	if err != nil && err != memcache.ErrCacheMiss {
		return err
	}

	if item != nil { // Exist
		err := p.config.UnSerializeFunc(store.DataPointer(), item.Value)
		if err != nil {
			return err
		}
	}

	releaseItem(item)

	return nil
}

func (p *Provider) Save(store *session.Store) error {
	expiration := int32(store.GetExpiration().Seconds())
	if expiration > math.MaxInt32 {
		return errExpirationIsTooBig
	}

	data := store.GetAll()
	value, err := p.config.SerializeFunc(data)
	if err != nil {
		return err
	}

	item := acquireItem()
	item.Key = p.getMemCacheSessionKey(store.GetSessionID())
	item.Value = value
	item.Expiration = expiration

	err = p.db.Set(item)

	releaseItem(item)

	return nil
}

// Regenerate regenerate session
func (p *Provider) Regenerate(id []byte, newStore *session.Store) error {
	key := p.getMemCacheSessionKey(id)
	newKey := p.getMemCacheSessionKey(newStore.GetSessionID())

	item, err := p.db.Get(key)
	if err != nil && err != memcache.ErrCacheMiss {
		return err
	}

	if item != nil { // Exist
		newItem := acquireItem()
		newItem.Key = newKey
		newItem.Value = item.Value
		newItem.Expiration = item.Expiration

		if err = p.db.Set(newItem); err != nil {
			return err
		}

		if err = p.db.Delete(key); err != nil {
			return err
		}

		if err := p.config.UnSerializeFunc(newStore.DataPointer(), newItem.Value); err != nil {
			return err
		}

		releaseItem(newItem)
	}

	releaseItem(item)

	return nil
}

// Destroy destroy session by sessionID
func (p *Provider) Destroy(id []byte) error {
	key := p.getMemCacheSessionKey(id)
	return p.db.Delete(key)
}

// Count session values count
func (p *Provider) Count() int {
	return 0
}

// NeedGC not need gc
func (p *Provider) NeedGC() bool {
	return false
}

// GC session memcache provider not need garbage collection
func (p *Provider) GC() {}
