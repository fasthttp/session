package postgres

import (
	"sync"

	"github.com/fasthttp/session"
)

// Config session postgres configuration
type Config struct {

	// The host to connect to. Values that start with / are for unix domain sockets. (default is localhost)
	Host string

	// The port to bind to. (default is 5432)
	Port int64

	// postgres user to sign in as
	Username string

	// postgres user's password
	Password string

	// Maximum wait for connection, in seconds. Zero or
	// not specified means wait indefinitely.
	ConnTimeout int64

	// name of the database to connect to
	Database string

	// session table name
	TableName string

	// postgres max free idle
	SetMaxIdleConn int

	// postgres max open idle
	SetMaxOpenConn int

	// session value serialize func
	SerializeFunc func(src session.Dict) ([]byte, error)

	// session value unSerialize func
	UnSerializeFunc func(dst *session.Dict, src []byte) error
}

// Provider provider struct
type Provider struct {
	config      *Config
	db          *Dao
	maxLifeTime int64

	storePool sync.Pool
}

// Store store struct
type Store struct {
	session.Store
}

// Dao database access object
type Dao struct {
	session.Dao

	tableName string

	sqlGetSessionBySessionID      string
	sqlCountSessions              string
	sqlUpdateBySessionID          string
	sqlDeleteBySessionID          string
	sqlDeleteSessionByMaxLifeTime string
	sqlInsert                     string
	sqlRegenerate                 string
}

// DBRow database row definition
type DBRow struct {
	sessionID  string
	contents   string
	lastActive int
}
