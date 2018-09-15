package memory

import (
	"errors"
	"reflect"
	"time"

	"github.com/fasthttp/session"
	"github.com/savsgio/dictpool"
)

// session memory provider

// ProviderName memory provider name
const ProviderName = "memory"

// Provider provider struct
type Provider struct {
	config      *Config
	values      *dictpool.Dict
	maxLifeTime int64
}

// NewProvider new memory provider
func NewProvider() *Provider {
	return &Provider{
		config:      &Config{},
		values:      dictpool.AcquireDict(),
		maxLifeTime: 0,
	}
}

// Init init provider config
func (mp *Provider) Init(lifeTime int64, memoryConfig session.ProviderConfig) error {
	if memoryConfig.Name() != ProviderName {
		return errors.New("session memory provider init error, config must memory config")
	}
	vc := reflect.ValueOf(memoryConfig)
	mc := vc.Interface().(*Config)
	mp.config = mc

	mp.maxLifeTime = lifeTime
	return nil
}

// NeedGC need gc
func (mp *Provider) NeedGC() bool {
	return true
}

// GC session garbage collection
func (mp *Provider) GC() {
	for _, kv := range mp.values.D {
		if time.Now().Unix() >= kv.Value.(*Store).lastActiveTime+mp.maxLifeTime {
			// destroy session sessionID
			mp.Destroy(string(kv.Key))
			return
		}
	}
}

// ReadStore read session store by session id
func (mp *Provider) ReadStore(sessionID string) (session.SessionStore, error) {
	memStore := mp.values.Get(sessionID)
	if memStore != nil {
		return memStore.(*Store), nil
	}

	newMemStore := NewMemoryStore(sessionID)
	mp.values.Set(sessionID, newMemStore)

	return newMemStore, nil
}

// Regenerate regenerate session
func (mp *Provider) Regenerate(oldSessionID string, sessionID string) (session.SessionStore, error) {
	memStoreInter := mp.values.Get(oldSessionID)
	if memStoreInter != nil {
		memStore := memStoreInter.(*Store)
		// insert new session store
		newMemStore := NewMemoryStoreData(sessionID, memStore.GetAll())
		mp.values.Set(sessionID, newMemStore)
		// delete old session store
		mp.values.Del(oldSessionID)

		return newMemStore, nil
	}

	memStore := NewMemoryStore(sessionID)
	mp.values.Set(sessionID, memStore)

	return memStore, nil
}

// Destroy destroy session by sessionID
func (mp *Provider) Destroy(sessionID string) error {
	mp.values.Del(sessionID)
	return nil
}

// Count session values count
func (mp *Provider) Count() int {
	return len(mp.values.D)
}

// register session provider
func init() {
	session.Register(ProviderName, NewProvider())
}
