package mysql

import (
	"fmt"
	"net/url"
)

// Config session mysql config
type Config struct {

	// mysql server host
	Host string

	// mysql server port
	Port int

	// mysql username
	Username string

	// mysql password
	Password string

	// mysql conn charset
	Charset string

	// mysql Collate
	Collate string

	// database name
	Database string

	// session table name
	TableName string

	// mysql conn timeout(s)
	Timeout int

	// mysql read timeout(s)
	ReadTimeout int

	// mysql write timeout(s)
	WriteTimeout int

	// mysql max free idle
	SetMaxIdleConn int

	// mysql max open idle
	SetMaxOpenConn int

	// session value serialize func
	SerializeFunc func(data map[string]interface{}) ([]byte, error)

	// session value unSerialize func
	UnSerializeFunc func(data []byte) (map[string]interface{}, error)
}

// NewConfigWith instance new config with especific paremters
func NewConfigWith(host string, port int, user, pass, dbName, tableName string) (cf *Config) {
	cf = NewDefaultConfig()
	cf.Host = host
	cf.Port = port
	cf.Username = user
	cf.Password = pass
	cf.Database = dbName
	cf.TableName = tableName
	return
}

// NewDefaultConfig return default config instance
func NewDefaultConfig() *Config {
	return &Config{
		Charset:        "utf8",
		Collate:        "utf8_general_ci",
		Database:       "test",
		TableName:      "test",
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

// Name return provider name
func (mc *Config) Name() string {
	return ProviderName
}
