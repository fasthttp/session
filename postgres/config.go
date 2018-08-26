package postgres

import (
	"fmt"
	"net/url"
)

// Config session postgres config
type Config struct {

	// The host to connect to. Values that start with / are for unix domain sockets. (default is localhost)
	Host string

	// The port to bind to. (default is 5432)
	Port int64

	// postgres user to sign in as
	Username string

	// postgres user's password
	Password string

	// Maximum wait for connection, in seconds. Zero or
	// not specified means wait indefinitely.
	ConnTimeout int64

	// name of the database to connect to
	Database string

	// session table name
	TableName string

	// postgres max free idle
	SetMaxIdleConn int

	// postgres max open idle
	SetMaxOpenConn int

	// session value serialize func
	SerializeFunc func(data map[string]interface{}) ([]byte, error)

	// session value unSerialize func
	UnSerializeFunc func(data []byte) (map[string]interface{}, error)
}

// NewConfigWith instance new config with especific paremters
func NewConfigWith(host string, port int64, username string, password string, dbName string, tableName string) (cf *Config) {
	cf = NewDefaultConfig()
	cf.Host = host
	cf.Port = port
	cf.Username = username
	cf.Password = password
	cf.Database = dbName
	cf.TableName = tableName
	return
}

// NewDefaultConfig return default config instance
func NewDefaultConfig() *Config {
	return &Config{
		Host:           "127.0.0.1",
		Port:           5432,
		Username:       "root",
		Password:       "",
		ConnTimeout:    3000,
		Database:       "test",
		TableName:      "test",
		SetMaxOpenConn: 500,
		SetMaxIdleConn: 50,
	}
}

func (pc *Config) getPostgresDSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?connect_timeout=%d&sslmode=disable",
		url.QueryEscape(pc.Username),
		pc.Password,
		url.QueryEscape(pc.Host),
		pc.Port,
		url.QueryEscape(pc.Database),
		pc.ConnTimeout)
}

// Name return provider name
func (pc *Config) Name() string {
	return ProviderName
}
