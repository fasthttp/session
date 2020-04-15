package postgre

import (
	"fmt"
	"net/url"
)

func (pc *Config) getPostgresDSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?connect_timeout=%d&sslmode=disable",
		url.QueryEscape(pc.Username),
		pc.Password,
		url.QueryEscape(pc.Host),
		pc.Port,
		url.QueryEscape(pc.Database),
		pc.ConnTimeout)
}

// NewConfigWith instance new configuration with especific paremters
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

// NewDefaultConfig return default configuration
func NewDefaultConfig() Config {
	return Config{
		Host:           "127.0.0.1",
		Port:           5432,
		Username:       "root",
		Password:       "",
		ConnTimeout:    3000,
		Database:       "session",
		TableName:      "session",
		SetMaxOpenConn: 500,
		SetMaxIdleConn: 50,
	}
}
