package redis

import (
	"sync"
	"time"

	"github.com/fasthttp/session"
	"github.com/go-redis/redis"
)

// Config session redis config
type Config struct {

	// Redis server host
	Host string

	// Redis server port
	Port int64

	// Maximum number of socket connections.
	PoolSize int

	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	// (s)
	IdleTimeout int64

	// redis server conn auth, default ""
	Password string

	// select db number, default 0
	DbNumber int

	// sessionID as redis key prefix
	KeyPrefix string

	// session value serialize func
	SerializeFunc func(src session.Dict) ([]byte, error)

	// session value unSerialize func
	UnSerializeFunc func(dst *session.Dict, src []byte) error
}

// Provider provider struct
type Provider struct {
	config      *Config
	db          *redis.Client
	maxLifeTime time.Duration

	storePool sync.Pool
}

// Store store struct
type Store struct {
	session.Store
}
