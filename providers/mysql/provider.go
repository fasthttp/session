package mysql

import (
	"time"

	"github.com/fasthttp/session/v2"
	"github.com/savsgio/gotils"
)

// New returns a new mysql provider configured
func New(cfg Config) (*Provider, error) {
	if cfg.Host == "" {
		return nil, errConfigHostEmpty
	}
	if cfg.Port == 0 {
		return nil, errConfigPortZero
	}

	if cfg.SerializeFunc == nil {
		cfg.SerializeFunc = session.Base64Encode
	}
	if cfg.UnSerializeFunc == nil {
		cfg.UnSerializeFunc = session.Base64Decode
	}

	db, err := newDao(cfg.getMysqlDSN(), cfg.TableName)
	if err != nil {
		return nil, err
	}

	db.Connection.SetMaxOpenConns(cfg.SetMaxIdleConn)
	db.Connection.SetMaxIdleConns(cfg.SetMaxIdleConn)

	if err := db.Connection.Ping(); err != nil {
		return nil, err
	}

	p := &Provider{
		config: cfg,
		db:     db,
	}

	return p, nil
}

// Get sets the user session to the given store
func (p *Provider) Get(store *session.Store) error {
	row, err := p.db.getSessionBySessionID(store.GetSessionID())
	if err != nil {
		return err
	}

	if row.sessionID != "" { // Exist
		err = p.config.UnSerializeFunc(store.Ptr(), gotils.S2B(row.contents))
		if err != nil {
			return err
		}

	} else { // Not exist
		_, err = p.db.insert(store.GetSessionID(), nil, time.Now().Unix(), store.GetExpiration())
		if err != nil {
			return err
		}
	}

	releaseDBRow(row)

	return nil
}

// Save saves the user session from the given store
func (p *Provider) Save(store *session.Store) error {
	data := store.GetAll()
	value, err := p.config.SerializeFunc(data)
	if err != nil {
		return err
	}

	_, err = p.db.updateBySessionID(store.GetSessionID(), value, time.Now().Unix(), store.GetExpiration())

	return err
}

// Regenerate updates a user session with the new session id
// and sets the user session to the store
func (p *Provider) Regenerate(id []byte, newStore *session.Store) error {
	row, err := p.db.getSessionBySessionID(id)
	if err != nil {
		return err
	}

	now := time.Now().Unix()

	if row.sessionID != "" { // Exists
		_, err = p.db.regenerate(id, newStore.GetSessionID(), now, newStore.GetExpiration())
		if err != nil {
			return err
		}

		err = p.config.UnSerializeFunc(newStore.Ptr(), gotils.S2B(row.contents))
		if err != nil {
			return err
		}

	} else { // Not exist
		_, err = p.db.insert(newStore.GetSessionID(), nil, now, newStore.GetExpiration())
		if err != nil {
			return err
		}
	}

	releaseDBRow(row)

	return nil
}

// Destroy destroys the user session from the given id
func (p *Provider) Destroy(id []byte) error {
	_, err := p.db.deleteBySessionID(id)
	return err
}

// Count returns the total of users sessions stored
func (p *Provider) Count() int {
	return p.db.countSessions()
}

// NeedGC indicates if the GC needs to be run
func (p *Provider) NeedGC() bool {
	return true
}

// GC destroys the expired user sessions
func (p *Provider) GC() {
	_, err := p.db.deleteExpiredSessions()
	if err != nil {
		panic(err)
	}
}
