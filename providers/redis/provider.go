package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/valyala/bytebufferpool"
)

var all = []byte("*")

// New returns a new configured redis provider
func New(cfg Config) (*Provider, error) {
	if cfg.Addr == "" {
		return nil, errConfigAddrEmpty
	}

	db := redis.NewClient(&redis.Options{
		Network:            cfg.Network,
		Addr:               cfg.Addr,
		Password:           cfg.Password,
		DB:                 cfg.DB,
		MaxRetries:         cfg.MaxRetries,
		MinRetryBackoff:    cfg.MinRetryBackoff,
		MaxRetryBackoff:    cfg.MaxRetryBackoff,
		DialTimeout:        cfg.DialTimeout,
		ReadTimeout:        cfg.ReadTimeout,
		WriteTimeout:       cfg.WriteTimeout,
		PoolSize:           cfg.PoolSize,
		MinIdleConns:       cfg.MinIdleConns,
		MaxConnAge:         cfg.MaxConnAge,
		PoolTimeout:        cfg.PoolTimeout,
		IdleTimeout:        cfg.IdleTimeout,
		IdleCheckFrequency: cfg.IdleCheckFrequency,
		TLSConfig:          cfg.TLSConfig,
		Limiter:            cfg.Limiter,
	})

	if err := db.Ping(context.Background()).Err(); err != nil {
		return nil, errRedisConnection(err)
	}

	p := &Provider{
		config: cfg,
		db:     db,
	}

	return p, nil
}

func (p *Provider) getRedisSessionKey(sessionID []byte) string {
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
	key := p.getRedisSessionKey(id)

	reply, err := p.db.Get(context.Background(), key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	return reply, nil

}

// Save saves the session data and expiration from the given session id
func (p *Provider) Save(id, data []byte, expiration time.Duration) error {
	key := p.getRedisSessionKey(id)

	return p.db.Set(context.Background(), key, data, expiration).Err()
}

// Regenerate updates the session id and expiration with the new session id
// of the the given current session id
func (p *Provider) Regenerate(id, newID []byte, expiration time.Duration) error {
	key := p.getRedisSessionKey(id)
	newKey := p.getRedisSessionKey(newID)

	exists, err := p.db.Exists(context.Background(), key).Result()
	if err != nil {
		return err
	}

	if exists > 0 { // Exist
		if err = p.db.Rename(context.Background(), key, newKey).Err(); err != nil {
			return err
		}

		if err = p.db.Expire(context.Background(), newKey, expiration).Err(); err != nil {
			return err
		}
	}

	return nil
}

// Destroy destroys the session from the given id
func (p *Provider) Destroy(id []byte) error {
	key := p.getRedisSessionKey(id)

	return p.db.Del(context.Background(), key).Err()
}

// Count returns the total of stored sessions
func (p *Provider) Count() int {
	reply, err := p.db.Keys(context.Background(), p.getRedisSessionKey(all)).Result()
	if err != nil {
		return 0
	}

	return len(reply)
}

// NeedGC indicates if the GC needs to be run
func (p *Provider) NeedGC() bool {
	return false
}

// GC destroys the expired sessions
func (p *Provider) GC() error {
	return nil
}
