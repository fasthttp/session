package memcache

import (
	"math"
	"sync"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
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

	db := memcache.New(cfg.ServerList...)
	db.Timeout = cfg.Timeout
	db.MaxIdleConns = cfg.MaxIdleConns

	if err := db.Ping(); err != nil {
		return nil, err
	}

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

// Get returns the data of the given session id
func (p *Provider) Get(id []byte) ([]byte, error) {
	key := p.getMemCacheSessionKey(id)

	item, err := p.db.Get(key)
	if err == memcache.ErrCacheMiss {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return item.Value, nil
}

// Save saves the session data and expiration from the given session id
func (p *Provider) Save(id, data []byte, expiration time.Duration) error {
	mcExpiration := int32(expiration.Seconds())
	if mcExpiration > math.MaxInt32 {
		return errExpirationIsTooBig
	}

	item := acquireItem()
	item.Key = p.getMemCacheSessionKey(id)
	item.Value = data
	item.Expiration = mcExpiration

	err := p.db.Set(item)

	releaseItem(item)

	return err
}

// Regenerate updates the session id and expiration with the new session id
// of the the given current session id
func (p *Provider) Regenerate(id, newID []byte, expiration time.Duration) error {
	key := p.getMemCacheSessionKey(id)
	newKey := p.getMemCacheSessionKey(newID)

	item, err := p.db.Get(key)
	if err != nil && err != memcache.ErrCacheMiss {
		return err
	}

	if item != nil { // Exist
		mcExpiration := int32(expiration.Seconds())
		if mcExpiration > math.MaxInt32 {
			return errExpirationIsTooBig
		}

		newItem := acquireItem()
		newItem.Key = newKey
		newItem.Value = item.Value
		newItem.Expiration = mcExpiration

		if err = p.db.Set(newItem); err != nil {
			return err
		}

		if err = p.db.Delete(key); err != nil {
			return err
		}

		releaseItem(newItem)
	}

	return nil
}

// Destroy destroys the session from the given id
func (p *Provider) Destroy(id []byte) error {
	key := p.getMemCacheSessionKey(id)
	return p.db.Delete(key)
}

// Count returns the total of stored sessions
func (p *Provider) Count() int {
	return 0
}

// NeedGC indicates if the GC needs to be run
func (p *Provider) NeedGC() bool {
	return false
}

// GC destroys the expired sessions
func (p *Provider) GC() error {
	return nil
}
