package file

// Config session file config
type Config struct {

	// session file save path
	SavePath string

	// session file suffix
	Suffix string

	// session value serialize func
	SerializeFunc func(data map[string]interface{}) ([]byte, error)

	// session value unSerialize func
	UnSerializeFunc func(data []byte) (map[string]interface{}, error)
}

// Name return provider name
func (fc *Config) Name() string {
	return ProviderName
}
