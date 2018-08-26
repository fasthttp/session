package memcache

import (
	"errors"
	"reflect"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/fasthttp/session"
)

// session MemCache provider

// ProviderName memcache provider name
const ProviderName = "memcache"

var (
	provider = NewProvider()
	encrypt  = session.NewEncrypt()
)

// Provider provider struct
type Provider struct {
	config         *Config
	values         *session.CCMap
	memCacheClient *memcache.Client
	maxLifeTime    int64
}

// NewProvider new memcache provider
func NewProvider() *Provider {
	return &Provider{
		config:         &Config{},
		values:         session.NewDefaultCCMap(),
		memCacheClient: &memcache.Client{},
	}
}

// Init init provider config
func (mcp *Provider) Init(lifeTime int64, memCacheConfig session.ProviderConfig) error {
	if memCacheConfig.Name() != ProviderName {
		return errors.New("session memcache provider init error, config must memcache config")
	}
	vc := reflect.ValueOf(memCacheConfig)
	rc := vc.Interface().(*Config)
	mcp.config = rc

	// config check
	if len(mcp.config.ServerList) == 0 {
		return errors.New("session memcache provider init error, config ServerList not empty")
	}
	if mcp.config.MaxIdle <= 0 {
		return errors.New("session memcache provider init error, config MaxIdle must be more than 0")
	}
	// init config serialize func
	if mcp.config.SerializeFunc == nil {
		mcp.config.SerializeFunc = encrypt.GOBEncode
	}
	if mcp.config.UnSerializeFunc == nil {
		mcp.config.UnSerializeFunc = encrypt.GOBDecode
	}
	// create memcache client
	mcp.memCacheClient = memcache.New(mcp.config.ServerList...)
	mcp.memCacheClient.MaxIdleConns = mcp.config.MaxIdle
	mcp.maxLifeTime = lifeTime
	return nil
}

// NeedGC not need gc
func (mcp *Provider) NeedGC() bool {
	return false
}

// GC session memcache provider not need garbage collection
func (mcp *Provider) GC() {}

// ReadStore read session store by session id
func (mcp *Provider) ReadStore(sessionID string) (session.SessionStore, error) {

	memClient := mcp.getMemCacheClient()

	item, err := memClient.Get(mcp.getMemCacheSessionKey(sessionID))
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return NewMemCacheStore(sessionID), nil
		}

		return nil, err
	}
	if len(item.Value) == 0 {
		return NewMemCacheStore(sessionID), nil
	}

	data, err := mcp.config.UnSerializeFunc(item.Value)
	if err != nil {
		return nil, err
	}

	return NewMemCacheStoreData(sessionID, data), nil
}

// Regenerate regenerate session
func (mcp *Provider) Regenerate(oldSessionID string, sessionID string) (session.SessionStore, error) {

	memClient := mcp.getMemCacheClient()

	item, err := memClient.Get(mcp.getMemCacheSessionKey(oldSessionID))
	if err != nil || len(item.Value) == 0 {
		// false, old sessionID not exists
		err := memClient.Set(&memcache.Item{
			Key:        mcp.getMemCacheSessionKey(sessionID),
			Value:      []byte(""),
			Expiration: int32(mcp.maxLifeTime),
		})
		if err != nil {
			return nil, err
		}
		return NewMemCacheStore(sessionID), nil
	}
	// true, old sessionID exists, delete old sessionID
	err = memClient.Delete(mcp.getMemCacheSessionKey(oldSessionID))
	if err != nil {
		return nil, err
	}
	item.Key = mcp.getMemCacheSessionKey(sessionID)
	item.Expiration = int32(mcp.maxLifeTime)
	err = memClient.Set(item)
	if err != nil {
		return nil, err
	}

	return mcp.ReadStore(sessionID)
}

// Destroy destroy session by sessionID
func (mcp *Provider) Destroy(sessionID string) error {
	memClient := mcp.getMemCacheClient()
	return memClient.Delete(mcp.getMemCacheSessionKey(sessionID))
}

// Count session values count
func (mcp *Provider) Count() int {
	return 0
}

// get memcache session key, prefix:sessionID
func (mcp *Provider) getMemCacheSessionKey(sessionID string) string {
	return mcp.config.KeyPrefix + ":" + sessionID
}

func (mcp *Provider) getMemCacheClient() *memcache.Client {
	if mcp.memCacheClient == nil {
		mcp.memCacheClient = memcache.New(mcp.config.ServerList...)
	}
	return mcp.memCacheClient
}

// register session provider
func init() {
	session.Register(ProviderName, provider)
}
