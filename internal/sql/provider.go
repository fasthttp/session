package sql

import (
	"database/sql"
	"time"

	"github.com/savsgio/gotils"
)

// NewProvider returns a new configured sql provider
func NewProvider(cfg ProviderConfig) (*Provider, error) {
	db, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	p := &Provider{
		config: cfg,
		db:     db,
	}

	return p, nil
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
//
// Returns the number of rows affected by an update, insert, or delete.
// Not every database or database driver may support this.
func (p *Provider) Exec(query string, args ...interface{}) (int64, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()

		return 0, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()

		return 0, err
	}

	return result.RowsAffected()
}

// Close closes the database and prevents new queries from starting.
// Close then waits for all queries that have started processing on the server
// to finish.
//
// It is rare to Close a DB, as the DB handle is meant to be
// long-lived and shared between many goroutines.
func (p *Provider) Close() error {
	return p.db.Close()
}

// Get returns the data of the given session id
func (p *Provider) Get(id []byte) ([]byte, error) {
	result := p.db.QueryRow(p.config.SQLGet, gotils.B2S(id))

	data := []byte("")

	err := result.Scan(&data)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return data, nil
}

// Save saves the session data and expiration from the given session id
func (p *Provider) Save(id, data []byte, expiration time.Duration) error {
	now := time.Now().UnixNano()

	n, err := p.Exec(p.config.SQLSave, gotils.B2S(data), now, expiration.Seconds(), gotils.B2S(id))
	if err != nil {
		return err
	}

	if n == 0 { // Not exist
		_, err = p.Exec(p.config.SQLInsert, gotils.B2S(id), gotils.B2S(data), now, expiration.Seconds())
		if err != nil {
			return err
		}
	}

	return nil
}

// Regenerate updates the session id and expiration with the new session id
// of the the given current session id
func (p *Provider) Regenerate(id, newID []byte, expiration time.Duration) error {
	now := time.Now().UnixNano()

	n, err := p.Exec(p.config.SQLRegenerate, gotils.B2S(newID), now, expiration.Seconds(), gotils.B2S(id))
	if err != nil {
		return err
	}

	if n == 0 { // Not exist
		_, err = p.Exec(p.config.SQLInsert, gotils.B2S(newID), "", now, expiration.Seconds())
		if err != nil {
			return err
		}
	}

	return nil
}

// Destroy destroys the session from the given id
func (p *Provider) Destroy(id []byte) error {
	_, err := p.Exec(p.config.SQLDestroy, id)
	return err
}

// Count returns the total of stored sessions
func (p *Provider) Count() int {
	row := p.db.QueryRow(p.config.SQLCount)

	total := 0
	if err := row.Scan(&total); err != nil {
		return 0
	}

	return total
}

// NeedGC indicates if the GC needs to be run
func (p *Provider) NeedGC() bool {
	return true
}

// GC destroys the expired sessions
func (p *Provider) GC() {
	_, err := p.Exec(p.config.SQLGC, time.Now().UnixNano())
	if err != nil {
		panic(err)
	}
}
