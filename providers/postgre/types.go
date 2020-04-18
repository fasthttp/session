package postgre

import (
	"time"

	"github.com/fasthttp/session/v2/internal/sql"
)

// Config configuration of provider
type Config struct {
	// The host to connect to. Values that start with / are for unix domain sockets. (default is localhost)
	Host string

	// The port to bind to. (default is 5432)
	Port int64

	// postgres user to sign in as
	Username string

	// postgres user's password
	Password string

	// name of the database to connect to
	Database string

	// session table name
	TableName string

	// Maximum wait for connection, in seconds. Zero or
	// not specified means wait indefinitely.
	Timeout time.Duration

	// postgre max free idle
	MaxIdleConn int

	// postgre max open idle
	MaxOpenConn int

	// postgre conn max open idle
	ConnMaxLifetime time.Duration
}

// Provider backend manager
type Provider struct {
	config Config

	*sql.Provider
}
