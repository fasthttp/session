package mysql

import (
	"time"

	"github.com/fasthttp/session/v2/internal/sql"
)

// Config configuration of provider
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

	// mysql conn timeout
	Timeout time.Duration

	// mysql read timeout
	ReadTimeout time.Duration

	// mysql write timeout
	WriteTimeout time.Duration

	// mysql max free idle
	MaxIdleConn int

	// mysql max open idle
	MaxOpenConn int

	// mysql conn max open idle
	ConnMaxLifetime time.Duration
}

// Provider backend manager
type Provider struct {
	config Config

	*sql.Provider
}
