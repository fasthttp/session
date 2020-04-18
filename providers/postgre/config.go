package postgre

import (
	"fmt"
	"net/url"
	"time"
)

// NewDefaultConfig returns a default configuration
func NewDefaultConfig() Config {
	return Config{
		Host:            "127.0.0.1",
		Port:            5432,
		Username:        "root",
		Password:        "",
		Database:        "session",
		TableName:       "session",
		Timeout:         30 * time.Second,
		MaxOpenConn:     100,
		MaxIdleConn:     100,
		ConnMaxLifetime: 1 * time.Second,
	}
}

// NewConfigWith returns a new configuration with especific paremters
func NewConfigWith(host string, port int64, username string, password string, dbName string, tableName string) Config {
	cf := NewDefaultConfig()
	cf.Host = host
	cf.Port = port
	cf.Username = username
	cf.Password = password
	cf.Database = dbName
	cf.TableName = tableName

	return cf
}

func (c *Config) dsn() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?connect_timeout=%d&sslmode=disable",
		url.QueryEscape(c.Username),
		c.Password,
		url.QueryEscape(c.Host),
		c.Port,
		url.QueryEscape(c.Database),
		int64(c.Timeout.Seconds()))
}
