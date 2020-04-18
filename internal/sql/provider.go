package sql

import (
	"database/sql"
	"time"

	"github.com/savsgio/gotils"
)

// New returns a new mysql provider configured
func NewProvider(cfg ProviderConfig) (*Provider, error) {
	db, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetMaxIdleConns(cfg.MaxIdleConn)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	p := &Provider{
		config: cfg,
		DB:     db,
	}

	return p, nil
}

func (p *Provider) Exec(query string, args ...interface{}) (int64, error) {
	result, err := p.DB.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// Get sets the user session to the given store
func (p *Provider) Get(id []byte) ([]byte, error) {
	result := p.QueryRow(p.config.SQLGet, gotils.B2S(id))

	data := []byte("")

	err := result.Scan(&data)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return data, nil
}

// Save saves the user session from the given store
func (p *Provider) Save(id, data []byte, expiration time.Duration) error {
	now := time.Now().Unix()

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

// Regenerate updates a user session with the new session id
// and sets the user session to the store
func (p *Provider) Regenerate(id, newID []byte, expiration time.Duration) error {
	now := time.Now().Unix()

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

// Destroy destroys the user session from the given id
func (p *Provider) Destroy(id []byte) error {
	_, err := p.Exec(p.config.SQLDestroy, id)
	return err
}

// Count returns the total of users sessions stored
func (p *Provider) Count() int {
	row := p.QueryRow(p.config.SQLCount)

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

// GC destroys the expired user sessions
func (p *Provider) GC() {
	_, err := p.Exec(p.config.SQLGC, time.Now().Unix())
	if err != nil {
		panic(err)
	}
}
