package postgre

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	// Import postgres driver
	_ "github.com/lib/pq"
	"github.com/savsgio/gotils"
)

var dbRowPool = &sync.Pool{
	New: func() interface{} {
		return new(dbRow)
	},
}

func acquireDBRow() *dbRow {
	return dbRowPool.Get().(*dbRow)
}

func releaseDBRow(row *dbRow) {
	row.reset()
	dbRowPool.Put(row)
}

func (row *dbRow) reset() {
	row.sessionID = ""
	row.contents = ""
	row.lastActive = 0
}

func newDao(dsn, tableName string) (*dao, error) {
	db := &dao{tableName: tableName}
	db.Driver = "postgres"
	db.Dsn = dsn

	var err error
	db.Connection, err = sql.Open(db.Driver, db.Dsn)

	db.sqlGetSessionBySessionID = fmt.Sprintf("SELECT session_id,contents,last_active,expiration FROM %s WHERE session_id=$1", tableName)
	db.sqlCountSessions = fmt.Sprintf("SELECT count(*) as total FROM %s", tableName)
	db.sqlUpdateBySessionID = fmt.Sprintf("UPDATE %s SET contents=$1,last_active=$2,expiration=$3 WHERE session_id=$4", tableName)
	db.sqlDeleteBySessionID = fmt.Sprintf("DELETE FROM %s WHERE session_id=$1", tableName)
	db.sqlDeleteExpiredSessions = fmt.Sprintf("DELETE FROM %s WHERE last_active+expiration<=$1 AND expiration<>0", tableName)
	db.sqlInsert = fmt.Sprintf("INSERT INTO %s (session_id, contents, last_active, expiration) VALUES ($1,$2,$3,$4)", tableName)
	db.sqlRegenerate = fmt.Sprintf("UPDATE %s SET session_id=$1,last_active=$2,expiration=$3 WHERE session_id=$4", tableName)

	return db, err
}

func (db *dao) getSessionBySessionID(sessionID []byte) (*dbRow, error) {
	data := acquireDBRow()

	row, err := db.QueryRow(db.sqlGetSessionBySessionID, gotils.B2S(sessionID))
	if err != nil {
		return nil, err
	}

	err = row.Scan(&data.sessionID, &data.contents, &data.lastActive, &data.expiration)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	data.expiration *= time.Second

	return data, nil
}

func (db *dao) countSessions() int {
	row, err := db.QueryRow(db.sqlCountSessions)
	if err != nil {
		return 0
	}

	var total int
	err = row.Scan(&total)
	if err != nil {
		return 0
	}

	return total
}

func (db *dao) updateBySessionID(sessionID, contents []byte, lastActiveTime int64, expiration time.Duration) (int64, error) {
	return db.Exec(db.sqlUpdateBySessionID, gotils.B2S(contents), lastActiveTime, expiration.Seconds(), gotils.B2S(sessionID))
}

func (db *dao) deleteBySessionID(sessionID []byte) (int64, error) {
	return db.Exec(db.sqlDeleteBySessionID, gotils.B2S(sessionID))
}

func (db *dao) deleteExpiredSessions() (int64, error) {
	return db.Exec(db.sqlDeleteExpiredSessions, time.Now().Unix())
}

func (db *dao) insert(sessionID, contents []byte, lastActiveTime int64, expiration time.Duration) (int64, error) {
	return db.Exec(db.sqlInsert, gotils.B2S(sessionID), gotils.B2S(contents), lastActiveTime, expiration.Seconds())
}

func (db *dao) regenerate(oldID, newID []byte, lastActiveTime int64, expiration time.Duration) (int64, error) {
	return db.Exec(db.sqlRegenerate, gotils.B2S(newID), lastActiveTime, expiration.Seconds(), gotils.B2S(oldID))
}
