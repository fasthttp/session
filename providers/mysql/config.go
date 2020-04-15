package mysql

import (
	"fmt"
	"net/url"
)

// get mysql dsn
func (mc *Config) getMysqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%dms&readTimeout=%dms&writeTimeout=%dms&charset=%s&collation=%s",
		url.QueryEscape(mc.Username),
		mc.Password,
		url.QueryEscape(mc.Host),
		mc.Port,
		url.QueryEscape(mc.Database),
		mc.Timeout,
		mc.ReadTimeout,
		mc.WriteTimeout,
		url.QueryEscape(mc.Charset),
		url.QueryEscape(mc.Collate))
}

// NewConfigWith return new configuration with especific paremters
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

// NewDefaultConfig return default configuration
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
