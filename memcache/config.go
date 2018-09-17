package memcache

import "github.com/savsgio/dictpool"

// Config session memcache config
type Config struct {

	// memcache server list
	ServerList []string

	// MaxIdleConns specifies the maximum number of idle connections that will
	// be maintained per address. If less than one, DefaultMaxIdleConns will be
	// used.
	//
	// Consider your expected traffic rates and latency carefully. This should
	// be set to a number higher than your peak parallel requests.
	MaxIdle int

	// sessionID as memcache key prefix
	KeyPrefix string

	// session value serialize func
	SerializeFunc func(data *dictpool.Dict) ([]byte, error)

	// session value unSerialize func
	UnSerializeFunc func(data []byte) (*dictpool.Dict, error)
}

// Name return provider name
func (mc *Config) Name() string {
	return ProviderName
}
