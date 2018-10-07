package sqlite3

import (
	"github.com/fasthttp/session"
)

// Config session sqlite3 configuration
type Config struct {

	// sqlite3 db file path
	DBPath string

	// session table name
	TableName string

	// sqlite3 max free idle
	SetMaxIdleConn int

	// sqlite3 max open idle
	SetMaxOpenConn int

	// session value serialize func
	SerializeFunc func(src *session.Dict) ([]byte, error)

	// session value unSerialize func
	UnSerializeFunc func(src []byte, dst *session.Dict) error
}

// Provider provider struct
type Provider struct {
	config      *Config
	db          *Dao
	maxLifeTime int64
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
