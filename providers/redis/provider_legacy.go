//go:build !go1.19
// +build !go1.19

package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var all = []byte("*")

// New returns a new configured redis provider
func New(cfg Config) (*Provider, error) {
	if cfg.Addr == "" {
		return nil, ErrConfigAddrEmpty
	}

	if cfg.Logger != nil {
		redis.SetLogger(cfg.Logger)
	}

	db := redis.NewClient(&redis.Options{
		Network:            cfg.Network,
		Addr:               cfg.Addr,
		Username:           cfg.Username,
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
		return nil, newErrRedisConnection(err)
	}

	p := &Provider{
		keyPrefix: cfg.KeyPrefix,
		db:        db,
	}

	return p, nil
}

// NewFailover returns a new redis provider using sentinel to determine the redis server to connect to.
func NewFailover(cfg FailoverConfig) (*Provider, error) {
	if cfg.MasterName == "" {
		return nil, ErrConfigMasterNameEmpty
	}

	if cfg.Logger != nil {
		redis.SetLogger(cfg.Logger)
	}

	db := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:         cfg.MasterName,
		SentinelAddrs:      cfg.SentinelAddrs,
		SentinelUsername:   cfg.SentinelUsername,
		SentinelPassword:   cfg.SentinelPassword,
		SlaveOnly:          cfg.SlaveOnly,
		Username:           cfg.Username,
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
	})

	if err := db.Ping(context.Background()).Err(); err != nil {
		return nil, newErrRedisConnection(err)
	}

	p := &Provider{
		keyPrefix: cfg.KeyPrefix,
		db:        db,
	}

	return p, nil
}

// NewFailoverCluster returns a new redis provider using a group of sentinels to determine the redis server to connect to.
func NewFailoverCluster(cfg FailoverConfig) (*Provider, error) {
	if cfg.MasterName == "" {
		return nil, ErrConfigMasterNameEmpty
	}

	if cfg.Logger != nil {
		redis.SetLogger(cfg.Logger)
	}

	db := redis.NewFailoverClusterClient(&redis.FailoverOptions{
		MasterName:         cfg.MasterName,
		SentinelAddrs:      cfg.SentinelAddrs,
		SentinelUsername:   cfg.SentinelUsername,
		SentinelPassword:   cfg.SentinelPassword,
		RouteByLatency:     cfg.RouteByLatency,
		RouteRandomly:      cfg.RouteRandomly,
		SlaveOnly:          cfg.SlaveOnly,
		Username:           cfg.Username,
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
	})

	if err := db.Ping(context.Background()).Err(); err != nil {
		return nil, newErrRedisConnection(err)
	}

	p := &Provider{
		keyPrefix: cfg.KeyPrefix,
		db:        db,
	}

	return p, nil
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
