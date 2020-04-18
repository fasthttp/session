package sqlite3

import (
	"time"

	"github.com/fasthttp/session/v2/internal/sql"
)

// Config configuration of provider
type Config struct {
	// sqlite3 db file path
	DBPath string

	// session table name
	TableName string

	// sqlite3 max free idle
	MaxIdleConn int

	// sqlite3 max open idle
	MaxOpenConn int

	// mysql conn max open idle
	ConnMaxLifetime time.Duration
}

// Provider backend manager
type Provider struct {
	config Config

	*sql.Provider
}
