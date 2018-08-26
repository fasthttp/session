package memory

// session memory config

// Config session memoryh config
type Config struct{}

// Name return provider name
func (mc *Config) Name() string {
	return ProviderName
}
