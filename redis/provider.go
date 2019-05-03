package redis

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/fasthttp/session"
	"github.com/go-redis/redis"
	"github.com/valyala/bytebufferpool"
)

var (
	provider = NewProvider()
	encrypt  = session.NewEncrypt()
	all      = []byte("*")
)

// NewProvider new redis provider
func NewProvider() *Provider {
	return &Provider{
		config: new(Config),
		db:     new(redis.Client),

		storePool: sync.Pool{
			New: func() interface{} {
				return new(Store)
			},
		},
	}
}

func (rp *Provider) acquireStore(sessionID []byte, expiration time.Duration) *Store {
	store := rp.storePool.Get().(*Store)
	store.Init(sessionID, expiration)

	return store
}

func (rp *Provider) releaseStore(store *Store) {
	store.Reset()
	rp.storePool.Put(store)
}

// Init init provider config
func (rp *Provider) Init(expiration time.Duration, cfg session.ProviderConfig) error {
	if cfg.Name() != ProviderName {
		return errors.New("session redis provider init error, config must redis config")
	}

	rp.config = cfg.(*Config)
	rp.expiration = expiration

	// config check
	if rp.config.Host == "" {
		return errConfigHostEmpty
	}
	if rp.config.Port == 0 {
		return errConfigPortZero
	}
	if rp.config.PoolSize <= 0 {
		return errConfigPoolSizeZero
	}
	if rp.config.IdleTimeout <= 0 {
		return errConfigIdleTimeoutZero
	}

	// init config serialize func
	if rp.config.SerializeFunc == nil {
		rp.config.SerializeFunc = encrypt.MSGPEncode
	}
	if rp.config.UnSerializeFunc == nil {
		rp.config.UnSerializeFunc = encrypt.MSGPDecode
	}

	// create redis conn pool
	rp.db = redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%d", rp.config.Host, rp.config.Port),
		Password:    rp.config.Password,
		DB:          rp.config.DbNumber,
		PoolSize:    rp.config.PoolSize,
		IdleTimeout: time.Duration(rp.config.IdleTimeout) * time.Second,
	})

	// check redis conn
	err := rp.db.Ping().Err()
	if err != nil {
		return errRedisConnection(err)
	}

	return nil
}

// get redis session key, prefix:sessionID
func (rp *Provider) getRedisSessionKey(sessionID []byte) string {
	key := bytebufferpool.Get()
	key.SetString(rp.config.KeyPrefix)
	key.WriteString(":")
	key.Write(sessionID)

	keyStr := key.String()

	bytebufferpool.Put(key)

	return keyStr
}

// Get read session store by session id
func (rp *Provider) Get(sessionID []byte) (session.Storer, error) {
	store := rp.acquireStore(sessionID, rp.expiration)
	key := rp.getRedisSessionKey(sessionID)

	reply, err := rp.db.Get(key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(reply) > 0 { // Exist
		err = rp.config.UnSerializeFunc(store.DataPointer(), reply)
		if err != nil {
			return nil, err
		}
	}

	return store, nil

}

// Put put store into the pool.
func (rp *Provider) Put(store session.Storer) {
	rp.releaseStore(store.(*Store))
}

// Regenerate regenerate session
func (rp *Provider) Regenerate(oldID, newID []byte) (session.Storer, error) {
	oldKey := rp.getRedisSessionKey(oldID)
	newKey := rp.getRedisSessionKey(newID)

	exists, err := rp.db.Exists(oldKey).Result()
	if err != nil {
		return nil, err
	}

	if exists > 0 { // Exist
		err = rp.db.Rename(oldKey, newKey).Err()
		if err != nil {
			return nil, err
		}
		err = rp.db.Expire(newKey, rp.expiration).Err()
		if err != nil {
			return nil, err
		}
	}

	return rp.Get(newID)
}

// Destroy destroy session by sessionID
func (rp *Provider) Destroy(sessionID []byte) error {
	key := rp.getRedisSessionKey(sessionID)
	return rp.db.Del(key).Err()
}

// Count session values count
func (rp *Provider) Count() int {
	reply, err := rp.db.Keys(rp.getRedisSessionKey(all)).Result()
	if err != nil {
		return 0
	}

	return len(reply)
}

// NeedGC not need gc
func (rp *Provider) NeedGC() bool {
	return false
}

// GC session redis provider not need garbage collection
func (rp *Provider) GC() {}

// register session provider
func init() {
	err := session.Register(ProviderName, provider)
	if err != nil {
		panic(err)
	}
}
