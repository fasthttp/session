package sqlite3

import (
	"time"

	"github.com/fasthttp/session"
	"github.com/savsgio/gotils"
)

// New new sqlite3 provider
func New(cfg Config) (*Provider, error) {
	if cfg.DBPath == "" {
		return nil, errConfigDBPathEmpty
	}
	if cfg.SerializeFunc == nil {
		cfg.SerializeFunc = session.Base64Encode
	}
	if cfg.UnSerializeFunc == nil {
		cfg.UnSerializeFunc = session.Base64Decode
	}

	db, err := NewDao(cfg.DBPath, cfg.TableName)
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

// Get read session store by session id
func (p *Provider) Get(store *session.Store) error {
	row, err := p.db.getSessionBySessionID(store.GetSessionID())
	if err != nil {
		return err
	}

	if row.sessionID != "" { // Exist
		err = p.config.UnSerializeFunc(store.DataPointer(), gotils.S2B(row.contents))
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

// Put put store into the pool.
func (p *Provider) Save(store *session.Store) error {
	data := store.GetAll()
	value, err := p.config.SerializeFunc(data)
	if err != nil {
		return err
	}

	_, err = p.db.updateBySessionID(store.GetSessionID(), value, time.Now().Unix(), store.GetExpiration())

	return err
}

// Regenerate regenerate session
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

		err = p.config.UnSerializeFunc(newStore.DataPointer(), gotils.S2B(row.contents))
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

// Destroy destroy session by sessionID
func (p *Provider) Destroy(id []byte) error {
	_, err := p.db.deleteBySessionID(id)
	return err
}

// Count session values count
func (p *Provider) Count() int {
	return p.db.countSessions()
}

// NeedGC need gc
func (p *Provider) NeedGC() bool {
	return true
}

// GC session garbage collection
func (p *Provider) GC() {
	_, err := p.db.deleteExpiredSessions()
	if err != nil {
		panic(err)
	}
}
