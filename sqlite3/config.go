package sqlite3

import "github.com/savsgio/dictpool"

// Config session sqlite3 config
type Config struct {

	// sqlite3 db file path
	DBPath string

	// session table name
	TableName string

	// sqlite3 max free idle
	SetMaxIdleConn int

	// sqlite3 max open idle
	SetMaxOpenConn int

	// session value serialize func
	SerializeFunc func(data *dictpool.Dict) ([]byte, error)

	// session value unSerialize func
	UnSerializeFunc func(data []byte) (*dictpool.Dict, error)
}

// NewConfigWith instance new config with especific paremters
func NewConfigWith(dbPath, tableName string) (cf *Config) {
	cf = &Config{
		SetMaxOpenConn: 500,
		SetMaxIdleConn: 50,
	}
	cf.DBPath = dbPath
	cf.TableName = tableName
	return
}

// Name return provider name
func (sc *Config) Name() string {
	return ProviderName
}
