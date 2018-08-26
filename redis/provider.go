package redis

import (
	"errors"
	"reflect"

	"github.com/fasthttp/session"
	"github.com/gomodule/redigo/redis"
)

// session redis provider

// ProviderName redis provider name
const ProviderName = "redis"

var (
	provider = NewProvider()
	encrypt  = session.NewEncrypt()
)

// Provider provider struct
type Provider struct {
	config      *Config
	redisPool   *redis.Pool
	maxLifeTime int64
}

// NewProvider new redis provider
func NewProvider() *Provider {
	return &Provider{
		config:    &Config{},
		redisPool: &redis.Pool{},
	}
}

// Init init provider config
func (rp *Provider) Init(lifeTime int64, redisConfig session.ProviderConfig) error {
	if redisConfig.Name() != ProviderName {
		return errors.New("session redis provider init error, config must redis config")
	}
	vc := reflect.ValueOf(redisConfig)
	rc := vc.Interface().(*Config)
	rp.config = rc
	rp.maxLifeTime = lifeTime

	// config check
	if rp.config.Host == "" {
		return errors.New("session redis provider init error, config Host not empty")
	}
	if rp.config.Port == 0 {
		return errors.New("session redis provider init error, config Port not empty")
	}
	if rp.config.MaxIdle <= 0 {
		return errors.New("session redis provider init error, config MaxIdle must be more than 0")
	}
	if rp.config.IdleTimeout <= 0 {
		return errors.New("session redis provider init error, config IdleTimeout must be more than 0")
	}
	// init config serialize func
	if rp.config.SerializeFunc == nil {
		rp.config.SerializeFunc = encrypt.GOBEncode
	}
	if rp.config.UnSerializeFunc == nil {
		rp.config.UnSerializeFunc = encrypt.GOBDecode
	}
	// create redis conn pool
	rp.redisPool = newRedisPool(rp.config)

	// check redis conn
	conn := rp.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	if err != nil {
		return errors.New("session redis provider init error, " + err.Error())
	}
	return nil
}

// NeedGC not need gc
func (rp *Provider) NeedGC() bool {
	return false
}

// GC session redis provider not need garbage collection
func (rp *Provider) GC() {}

// ReadStore read session store by session id
func (rp *Provider) ReadStore(sessionID string) (session.SessionStore, error) {

	conn := rp.redisPool.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", rp.getRedisSessionKey(sessionID)))
	if err != nil && err != redis.ErrNil {
		return nil, err
	}
	if len(reply) == 0 {
		conn.Do("SET", rp.getRedisSessionKey(sessionID), "", "EX", rp.maxLifeTime)
		return NewRedisStore(sessionID), nil
	}

	data, err := rp.config.UnSerializeFunc(reply)
	if err != nil {
		return nil, err
	}

	return NewRedisStoreData(sessionID, data), nil
}

// Regenerate regenerate session
func (rp *Provider) Regenerate(oldSessionId string, sessionID string) (session.SessionStore, error) {

	conn := rp.redisPool.Get()
	defer conn.Close()

	existed, err := redis.Int(conn.Do("EXISTS", rp.getRedisSessionKey(oldSessionId)))
	if err != nil || existed == 0 {
		// false
		conn.Do("SET", rp.getRedisSessionKey(sessionID), "", "EX", rp.maxLifeTime)
		return NewRedisStore(sessionID), nil
	}
	// true
	conn.Do("RENAME", rp.getRedisSessionKey(oldSessionId), rp.getRedisSessionKey(sessionID))
	conn.Do("EXPIRE", rp.getRedisSessionKey(sessionID), rp.maxLifeTime)

	return rp.ReadStore(sessionID)
}

// Destroy destroy session by sessionID
func (rp *Provider) Destroy(sessionID string) error {
	conn := rp.redisPool.Get()
	defer conn.Close()

	existed, err := redis.Int(conn.Do("EXISTS", rp.getRedisSessionKey(sessionID)))
	if err != nil || existed == 0 {
		return nil
	}
	conn.Do("DEL", rp.getRedisSessionKey(sessionID))
	return nil
}

// Count session values count
func (rp *Provider) Count() int {
	conn := rp.redisPool.Get()
	defer conn.Close()

	replyMap, err := redis.Strings(conn.Do("KEYS", rp.config.KeyPrefix+":*"))
	if err != nil {
		return 0
	}
	return len(replyMap)
}

// get redis session key, prefix:sessionID
func (rp *Provider) getRedisSessionKey(sessionID string) string {
	return rp.config.KeyPrefix + ":" + sessionID
}

// register session provider
func init() {
	session.Register(ProviderName, provider)
}
