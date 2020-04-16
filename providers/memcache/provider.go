package memcache

import (
	"math"
	"sync"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/fasthttp/session/v2"
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

// New returns a new memcache provider configured
func New(cfg Config) (*Provider, error) {
	if len(cfg.ServerList) == 0 {
		return nil, errConfigServerListEmpty
	}
	if cfg.MaxIdleConns <= 0 {
		return nil, errConfigMaxIdleConnsZero
	}

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

// Get sets the user session to the given store
func (p *Provider) Get(store *session.Store) error {
	key := p.getMemCacheSessionKey(store.GetSessionID())

	item, err := p.db.Get(key)
	if err != nil && err != memcache.ErrCacheMiss {
		return err
	}

	if item != nil { // Exist
		err := p.config.UnSerializeFunc(store.Ptr(), item.Value)
		if err != nil {
			return err
		}
	}

	releaseItem(item)

	return nil
}

// Save saves the user session from the given store
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

	return err
}

// Regenerate updates a user session with the new session id
// and sets the user session to the store
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

		if err := p.config.UnSerializeFunc(newStore.Ptr(), newItem.Value); err != nil {
			return err
		}

		releaseItem(newItem)
	}

	releaseItem(item)

	return nil
}

// Destroy destroys the user session from the given id
func (p *Provider) Destroy(id []byte) error {
	key := p.getMemCacheSessionKey(id)
	return p.db.Delete(key)
}

// Count returns the total of users sessions stored
func (p *Provider) Count() int {
	return 0
}

// NeedGC indicates if the GC needs to be run
func (p *Provider) NeedGC() bool {
	return false
}

// GC destroys the expired user sessions
func (p *Provider) GC() {}
