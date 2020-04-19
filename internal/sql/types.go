package sql

import (
	"database/sql"
	"time"
)

// ProviderConfig provider settings
type ProviderConfig struct {
	// SQL driver
	Driver string

	// DB connection string
	DSN string

	// The maximum number of connections in the idle connection pool.
	//
	// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns,
	// then the new MaxIdleConns will be reduced to match the MaxOpenConns limit.
	//
	// If n <= 0, no idle connections are retained.
	//
	// The default max idle connections is currently 2. This may change in
	// a future release.
	MaxIdleConns int

	// The maximum number of open connections to the database.
	//
	// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
	// MaxIdleConns, then MaxIdleConns will be reduced to match the new
	// MaxOpenConns limit.
	//
	// If n <= 0, then there is no limit on the number of open connections.
	// The default is 0 (unlimited).
	MaxOpenConns int

	// The maximum amount of time a connection may be reused.
	//
	// Expired connections may be closed lazily before reuse.
	//
	// If d <= 0, connections are reused forever.
	ConnMaxLifetime time.Duration

	SQLGet        string
	SQLSave       string
	SQLRegenerate string
	SQLDestroy    string
	SQLCount      string
	SQLInsert     string
	SQLGC         string
}

// Provider backend manager
type Provider struct {
	config ProviderConfig
	db     *sql.DB
}
