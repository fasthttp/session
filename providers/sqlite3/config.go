package sqlite3

// NewConfigWith returns a new configuration with especific paremters
func NewConfigWith(dbPath, tableName string) Config {
	cf := NewDefaultConfig()
	cf.DBPath = dbPath
	cf.TableName = tableName

	return cf
}

// NewDefaultConfig returns a default configuration
func NewDefaultConfig() Config {
	return Config{
		DBPath:         "./",
		TableName:      "session",
		SetMaxOpenConn: 500,
		SetMaxIdleConn: 50,
	}
}
