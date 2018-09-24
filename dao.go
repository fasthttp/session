package session

import (
	"database/sql"
)

// NewDao create new database access object
func NewDao(driver, dsn string) (*Dao, error) {
	db := &Dao{Driver: driver, Dsn: dsn}

	var err error
	db.Connection, err = sql.Open(db.Driver, db.Dsn)

	return db, err
}

func (db *Dao) makeTxStmt(query string) (*sql.Tx, *sql.Stmt, error) {
	tx, err := db.Connection.Begin()
	if err != nil {
		return nil, nil, err
	}

	stmt, err := tx.Prepare(query)

	return tx, stmt, err
}

func (db *Dao) makeStmt(query string) (*sql.Stmt, error) {
	return db.Connection.Prepare(query)
}

// Exec insert or update data from database
func (db *Dao) Exec(query string, args ...interface{}) (int64, error) {
	tx, stmt, err := db.makeTxStmt(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	x, err := res.RowsAffected()
	if err != nil {
		return 0, nil
	}
	err = tx.Commit()

	return x, err
}

// Query get data from database
func (db *Dao) Query(query string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := db.makeStmt(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	dataSet, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}

	return dataSet, err
}

// QueryRow get just one data from database
func (db *Dao) QueryRow(query string, args ...interface{}) (*sql.Row, error) {
	stmt, err := db.makeStmt(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryRow(args...), nil
}
