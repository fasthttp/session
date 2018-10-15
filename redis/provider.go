package redis

import (
	"errors"

	"github.com/fasthttp/session"
	"github.com/gomodule/redigo/redis"
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
		config:    new(Config),
		redisPool: new(redis.Pool),
	}
}

// Init init provider config
func (rp *Provider) Init(lifeTime int64, cfg session.ProviderConfig) error {
	if cfg.Name() != ProviderName {
		return errors.New("session redis provider init error, config must redis config")
	}

	rp.config = cfg.(*Config)
	rp.maxLifeTime = lifeTime

	// config check
	if rp.config.Host == "" {
		return errConfigHostEmpty
	}
	if rp.config.Port == 0 {
		return errConfigPortZero
	}
	if rp.config.MaxIdle <= 0 {
		return errConfigMaxIdleZero
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
	rp.redisPool = newRedisPool(rp.config)

	// check redis conn
	conn := rp.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
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
	conn := rp.redisPool.Get()
	defer conn.Close()

	store := NewStore(sessionID)
	key := rp.getRedisSessionKey(sessionID)

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err == nil { // Exist
		err := rp.config.UnSerializeFunc(reply, store.GetDataPointer())
		if err != nil {
			return nil, err
		}

	} else if err == redis.ErrNil { // Not exist
		conn.Do("SET", key, "", "EX", rp.maxLifeTime)
		err = nil
	}

	return store, err

}

// Regenerate regenerate session
func (rp *Provider) Regenerate(oldID, newID []byte) (session.Storer, error) {
	conn := rp.redisPool.Get()
	defer conn.Close()

	oldKey := rp.getRedisSessionKey(oldID)
	newKey := rp.getRedisSessionKey(newID)

	exists, err := redis.Bool(conn.Do("EXISTS", oldKey))
	if err != nil || !exists { // Not exist
		conn.Do("SET", newKey, "", "EX", rp.maxLifeTime)
		return NewStore(newID), nil
	}

	// Exist
	conn.Do("RENAME", oldKey, newKey)
	conn.Do("EXPIRE", newKey, rp.maxLifeTime)

	return rp.Get(newID)
}

// Destroy destroy session by sessionID
func (rp *Provider) Destroy(sessionID []byte) error {
	conn := rp.redisPool.Get()
	defer conn.Close()

	key := rp.getRedisSessionKey(sessionID)

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil || !exists {
		return nil
	}
	conn.Do("DEL", key)

	return nil
}

// Count session values count
func (rp *Provider) Count() int {
	conn := rp.redisPool.Get()
	defer conn.Close()

	reply, err := redis.ByteSlices(conn.Do("KEYS", rp.getRedisSessionKey(all)))
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
	session.Register(ProviderName, provider)
}
