package redis

import (
	"fmt"

	"github.com/fasthttp/session"
	"github.com/go-redis/redis/v7"
	"github.com/valyala/bytebufferpool"
)

var all = []byte("*")

// New new redis provider
func New(cfg Config) (*Provider, error) {
	// config check
	if cfg.Host == "" {
		return nil, errConfigHostEmpty
	}
	if cfg.Port == 0 {
		return nil, errConfigPortZero
	}
	if cfg.PoolSize <= 0 {
		return nil, errConfigPoolSizeZero
	}
	if cfg.IdleTimeout <= 0 {
		return nil, errConfigIdleTimeoutZero
	}

	// init config serialize func
	if cfg.SerializeFunc == nil {
		cfg.SerializeFunc = session.MSGPEncode
	}
	if cfg.UnSerializeFunc == nil {
		cfg.UnSerializeFunc = session.MSGPDecode
	}

	// create redis conn pool
	db := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:    cfg.Password,
		DB:          cfg.DbNumber,
		PoolSize:    cfg.PoolSize,
		IdleTimeout: cfg.IdleTimeout,
	})

	// check redis conn
	if err := db.Ping().Err(); err != nil {
		return nil, errRedisConnection(err)
	}

	p := &Provider{
		config: cfg,
		db:     db,
	}

	return p, nil
}

// get redis session key, prefix:sessionID
func (p *Provider) getRedisSessionKey(sessionID []byte) string {
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
	key := p.getRedisSessionKey(store.GetSessionID())

	reply, err := p.db.Get(key).Bytes()
	if err != nil && err != redis.Nil {
		return err
	}

	if len(reply) > 0 { // Exist
		err = p.config.UnSerializeFunc(store.DataPointer(), reply)
		if err != nil {
			return err
		}
	}

	return nil

}

// Put put store into the pool.
func (p *Provider) Save(store *session.Store) error {
	data := store.GetAll()
	b, err := p.config.SerializeFunc(data)
	if err != nil {
		return err
	}

	return p.db.Set(p.getRedisSessionKey(store.GetSessionID()), b, store.GetExpiration()).Err()
}

// Regenerate regenerate session
func (p *Provider) Regenerate(id []byte, newStore *session.Store) error {
	key := p.getRedisSessionKey(id)
	newKey := p.getRedisSessionKey(newStore.GetSessionID())

	exists, err := p.db.Exists(key).Result()
	if err != nil {
		return err
	}

	if exists > 0 { // Exist
		if err = p.db.Rename(key, newKey).Err(); err != nil {
			return err
		}

		if err = p.db.Expire(newKey, newStore.GetExpiration()).Err(); err != nil {
			return err
		}
	}

	return p.Get(newStore)
}

// Destroy destroy session by sessionID
func (p *Provider) Destroy(id []byte) error {
	key := p.getRedisSessionKey(id)
	return p.db.Del(key).Err()
}

// Count session values count
func (p *Provider) Count() int {
	reply, err := p.db.Keys(p.getRedisSessionKey(all)).Result()
	if err != nil {
		return 0
	}

	return len(reply)
}

// NeedGC not need gc
func (p *Provider) NeedGC() bool {
	return false
}

// GC session redis provider not need garbage collection
func (p *Provider) GC() {}
