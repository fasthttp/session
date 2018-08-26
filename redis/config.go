package redis

// Config session redis config
type Config struct {

	// Redis server host
	Host string

	// Redis server port
	Port int64

	// Maximum number of idle connections in the redis server pool.
	MaxIdle int

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
	SerializeFunc func(data map[string]interface{}) ([]byte, error)

	// session value unSerialize func
	UnSerializeFunc func(data []byte) (map[string]interface{}, error)
}

// Name return provider name
func (mc *Config) Name() string {
	return ProviderName
}
