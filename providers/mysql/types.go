package mysql

import (
	"time"

	"github.com/fasthttp/session/v2"
	gotilsDao "github.com/savsgio/gotils/dao"
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
