package mysql

import (
	"fmt"
	"net/url"
	"time"
)

// NewDefaultConfig returns a default configuration
func NewDefaultConfig() Config {
	return Config{
		Host:            "127.0.0.1",
		Port:            3306,
		Username:        "root",
		Password:        "",
		Database:        "session",
		TableName:       "session",
		DropTable:       false,
		Charset:         "utf8",
		Collation:       "utf8_general_ci",
		Timeout:         30 * time.Second,
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
		MaxOpenConns:    100,
		MaxIdleConns:    100,
		ConnMaxLifetime: 1 * time.Second,
	}
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

func (c *Config) dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%ds&readTimeout=%ds&writeTimeout=%ds&charset=%s&collation=%s",
		url.QueryEscape(c.Username),
		c.Password,
		url.QueryEscape(c.Host),
		c.Port,
		url.QueryEscape(c.Database),
		int64(c.Timeout.Seconds()),
		int64(c.ReadTimeout.Seconds()),
		int64(c.WriteTimeout.Seconds()),
		url.QueryEscape(c.Charset),
		url.QueryEscape(c.Collation))
}
