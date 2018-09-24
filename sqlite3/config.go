package sqlite3

// NewConfigWith instance new configuration with especific paremters
func NewConfigWith(dbPath, tableName string) *Config {
	cf := NewDefaultConfig()
	cf.DBPath = dbPath
	cf.TableName = tableName

	return cf
}

// NewDefaultConfig return default configuration
func NewDefaultConfig() *Config {
	cf := &Config{
		SetMaxOpenConn: 500,
		SetMaxIdleConn: 50,
	}

	return cf
}

// Name return provider name
func (sc *Config) Name() string {
	return ProviderName
}
