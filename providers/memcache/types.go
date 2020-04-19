package memcache

import (
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// Config configuration of provider
type Config struct {
	// Prefix key
	KeyPrefix string

	// Server list
	ServerList []string

	// The socket read/write timeout.
	// If zero, DefaultTimeout is used.
	Timeout time.Duration

	// The maximum number of idle connections that will
	// be maintained per address. If less than one, DefaultMaxIdleConns will be
	// used.
	//
	// Consider your expected traffic rates and latency carefully. This should
	// be set to a number higher than your peak parallel requests.
	MaxIdleConns int
}

// Provider backend manager
type Provider struct {
	config Config
	db     *memcache.Client
}
