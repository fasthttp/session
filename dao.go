package session

import (
	"database/sql"
	"sync"
)

var daoBufferPool = sync.Pool{
	New: func() interface{} {
		return new(daoBuffer)
	},
}

func (buf *daoBuffer) reset() {
	buf.tx = nil
	buf.stmt = nil
}

func acquireDaoBuffer() *daoBuffer {
	return daoBufferPool.Get().(*daoBuffer)
}

func releaseDaoBuffer(buf *daoBuffer) {
	buf.reset()
	daoBufferPool.Put(buf)
}

// NewDao create new database access object
func NewDao(driver, dsn string) (*Dao, error) {
	db := &Dao{Driver: driver, Dsn: dsn}

	var err error
	db.Connection, err = sql.Open(db.Driver, db.Dsn)

	return db, err
}

func (db *Dao) makeTxStmt(query string) (*daoBuffer, error) {
	var err error

	buf := acquireDaoBuffer()

	buf.tx, err = db.Connection.Begin()
	if err != nil {
		return buf, err
	}

	buf.stmt, err = buf.tx.Prepare(query)

	return buf, err
}

func (db *Dao) makeStmt(query string) (*daoBuffer, error) {
	var err error

	buf := acquireDaoBuffer()
	buf.stmt, err = db.Connection.Prepare(query)

	return buf, err
}

// Exec insert or update data from database
func (db *Dao) Exec(query string, args ...interface{}) (int64, error) {
	buf, err := db.makeTxStmt(query)
	if err != nil {
		return 0, err
	}
	defer buf.stmt.Close()

	res, err := buf.stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	x, err := res.RowsAffected()
	if err != nil {
		return 0, nil
	}
	err = buf.tx.Commit()

	releaseDaoBuffer(buf)

	return x, err
}

// Query get data from database
func (db *Dao) Query(query string, args ...interface{}) (*sql.Rows, error) {
	buf, err := db.makeStmt(query)
	if err != nil {
		return nil, err
	}
	defer buf.stmt.Close()

	rows, err := buf.stmt.Query(args...)
	if err != nil {
		return nil, err
	}

	releaseDaoBuffer(buf)

	return rows, err
}

// QueryRow get just one data from database
func (db *Dao) QueryRow(query string, args ...interface{}) (*sql.Row, error) {
	buf, err := db.makeStmt(query)
	if err != nil {
		return nil, err
	}
	defer buf.stmt.Close()

	return buf.stmt.QueryRow(args...), nil
}
