package sqlite3

import "time"

// NewDefaultConfig returns a default configuration
func NewDefaultConfig() Config {
	return Config{
		DBPath:          "./",
		TableName:       "session",
		MaxOpenConn:     100,
		MaxIdleConn:     100,
		ConnMaxLifetime: 1 * time.Second,
	}
}

// NewConfigWith returns a new configuration with especific paremters
func NewConfigWith(dbPath, tableName string) Config {
	cf := NewDefaultConfig()
	cf.DBPath = dbPath
	cf.TableName = tableName

	return cf
}
