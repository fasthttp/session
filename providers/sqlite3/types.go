package sqlite3

import (
	"time"

	"github.com/fasthttp/session/v2"
	gotilsDao "github.com/savsgio/gotils/dao"
)

// Config configuration of provider
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
