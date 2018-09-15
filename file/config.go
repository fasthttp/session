package file

import "github.com/savsgio/dictpool"

// Config session file config
type Config struct {

	// session file save path
	SavePath string

	// session file suffix
	Suffix string

	// session value serialize func
	SerializeFunc func(data *dictpool.Dict) ([]byte, error)

	// session value unSerialize func
	UnSerializeFunc func(data []byte) (*dictpool.Dict, error)
}

// Name return provider name
func (fc *Config) Name() string {
	return ProviderName
}
