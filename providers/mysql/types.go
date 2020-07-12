package mysql

import (
	"time"

	"github.com/fasthttp/session/v2/internal/sql"
)

// Config provider settings
type Config struct {
	// DB host
	Host string

	// DB port
	Port int

	// DB user name
	Username string

	// DB user password
	Password string

	// DB name
	Database string

	// DB table name
	TableName string

	// When set to true, this will Drop any existing table with the same name
	DropTable bool

	// Charset used for client-server interaction ("SET NAMES <value>")
	// If multiple charsets are set (separated by a comma),
	// the following charset is used if setting the charset failes
	//
	// This enables for example support for utf8mb4 (introduced in MySQL 5.5.3)
	// with fallback to utf8 for older servers (charset=utf8mb4,utf8).
	//
	// Usage of the charset parameter is discouraged because it issues additional queries to the server.
	// Unless you need the fallback behavior, please use collation instead.
	Charset string

	// Collation used for client-server interaction on connection
	// In contrast to charset, collation does not issue additional queries
	// If the specified collation is unavailable on the target server, the connection will fail.
	//
	// A list of valid charsets for a server is retrievable with SHOW COLLATION.
	//
	// The default collation (utf8mb4_general_ci) is supported from MySQL 5.5
	//You should use an older collation (e.g. utf8_general_ci) for older MySQL.
	//
	// Collations for charset "ucs2", "utf16", "utf16le", and "utf32" can not be used (ref).
	Collation string

	// Timeout for establishing connections, aka dial timeout, in seconds
	Timeout time.Duration

	// I/O read timeout, in seconds
	ReadTimeout time.Duration

	// I/O write timeout, in seconds
	WriteTimeout time.Duration

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
