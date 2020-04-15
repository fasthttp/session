package mysql

import (
	"fmt"
	"net/url"
)

func (c *Config) getMysqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%dms&readTimeout=%dms&writeTimeout=%dms&charset=%s&collation=%s",
		url.QueryEscape(c.Username),
		c.Password,
		url.QueryEscape(c.Host),
		c.Port,
		url.QueryEscape(c.Database),
		c.Timeout,
		c.ReadTimeout,
		c.WriteTimeout,
		url.QueryEscape(c.Charset),
		url.QueryEscape(c.Collate))
}

// NewConfigWith returns a new configuration with especific paremters
func NewConfigWith(host string, port int, user, pass, dbName, tableName string) Config {
	cf := NewDefaultConfig()
	cf.Host = host
	cf.Port = port
	cf.Username = user
	cf.Password = pass
	cf.Database = dbName
	cf.TableName = tableName

	return cf
}

// NewDefaultConfig returns a default configuration
func NewDefaultConfig() Config {
	return Config{
		Charset:        "utf8",
		Collate:        "utf8_general_ci",
		Database:       "session",
		TableName:      "session",
		Host:           "127.0.0.1",
		Port:           3306,
		Username:       "root",
		Password:       "",
		Timeout:        3000,
		ReadTimeout:    5000,
		WriteTimeout:   5000,
		SetMaxOpenConn: 500,
		SetMaxIdleConn: 50,
	}
}
