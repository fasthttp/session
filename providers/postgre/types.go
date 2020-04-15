package postgre

import (
	"time"

	"github.com/fasthttp/session"
	gotilsDao "github.com/savsgio/gotils/dao"
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

// Provider backend manager
type Provider struct {
	config Config
	db     *dao
}

type dao struct {
	gotilsDao.Dao

	tableName string

	sqlGetSessionBySessionID string
	sqlCountSessions         string
	sqlUpdateBySessionID     string
	sqlDeleteBySessionID     string
	sqlDeleteExpiredSessions string
	sqlInsert                string
	sqlRegenerate            string
}

type dbRow struct {
	sessionID  string
	contents   string
	lastActive int64
	expiration time.Duration
}
