package memory

import (
	"time"

	"github.com/fasthttp/session"
)

var provider = NewProvider()

// NewProvider new memory provider
func NewProvider() *Provider {
	return &Provider{
		config:      new(Config),
		values:      new(session.Dict),
		maxLifeTime: 0,
	}
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

// ReadStore read session by id
func (mp *Provider) ReadStore(sessionID []byte) (session.Storer, error) {
	currentStore := mp.values.GetBytes(sessionID)
	if currentStore != nil {
		return currentStore.(*Store), nil
	}

	newStore := NewStore(sessionID)
	mp.values.SetBytes(sessionID, newStore)

	return newStore, nil
}

// Regenerate regenerate session
func (mp *Provider) Regenerate(oldID, newID []byte) (session.Storer, error) {
	var store *Store

	val := mp.values.GetBytes(oldID)
	if val != nil {
		store = val.(*Store)
		store.SetSessionID(newID)
		mp.values.SetBytes(newID, store)
		mp.values.DelBytes(oldID)
	} else {
		store = NewStore(newID)
		mp.values.SetBytes(newID, store)
	}

	return store, nil
}

// Destroy destroy session by sessionID
func (mp *Provider) Destroy(sessionID []byte) error {
	val := mp.values.GetBytes(sessionID)
	if val != nil {
		releaseStore(val.(*Store))
	}

	mp.values.DelBytes(sessionID)

	return nil
}

// Count session values count
func (mp *Provider) Count() int {
	return len(mp.values.D)
}

// NeedGC need gc
func (mp *Provider) NeedGC() bool {
	return true
}

// GC session garbage collection
func (mp *Provider) GC() {
	for _, kv := range mp.values.D {
		if time.Now().Unix() >= (kv.Value.(*Store).lastActiveTime + mp.maxLifeTime) {
			mp.Destroy(kv.Key)
		}
	}
}

// register session provider
func init() {
	session.Register(ProviderName, provider)
}
