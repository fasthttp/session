package memcache

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/fasthttp/session"
)

// Config session memcache configuration
type Config struct {
	// ServerList memcache server list
	ServerList []string

	// MaxIdleConns specifies the maximum number of idle connections that will
	// be maintained per address. If less than one, DefaultMaxIdleConns will be
	// used.
	//
	// Consider your expected traffic rates and latency carefully. This should
	// be set to a number higher than your peak parallel requests.
	MaxIdleConns int

	// KeyPrefix sessionID as memcache key prefix
	KeyPrefix string

	// SerializeFunc session value serialize func
	SerializeFunc func(src *session.Dict) ([]byte, error)

	// UnSerializeFunc session value unSerialize func
	UnSerializeFunc func(src []byte, dst *session.Dict) error
}

// Provider provider struct
type Provider struct {
	config         *Config
	memCacheClient *memcache.Client
	maxLifeTime    int64
}

// Store store struct
type Store struct {
	session.Store
}
