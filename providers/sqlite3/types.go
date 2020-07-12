package sqlite3

import (
	"time"

	"github.com/fasthttp/session/v2/internal/sql"
)

// Config provider settings
type Config struct {
	// DB file path
	DBPath string

	// DB table name
	TableName string

	// When set to true, this will Drop any existing table with the same name
	DropTable bool

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
}

// Provider backend manager
type Provider struct {
	config Config

	*sql.Provider
}
