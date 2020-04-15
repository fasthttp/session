package mysql

import (
	"time"

	"github.com/fasthttp/session"
	"github.com/savsgio/gotils"
)

// New new mysql provider
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

	db, err := NewDao(cfg.getMysqlDSN(), cfg.TableName)
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
func (mp *Provider) Destroy(id []byte) error {
	_, err := mp.db.deleteBySessionID(id)
	return err
}

// Count session values count
func (mp *Provider) Count() int {
	return mp.db.countSessions()
}

// NeedGC need gc
func (mp *Provider) NeedGC() bool {
	return true
}

// GC session garbage collection
func (mp *Provider) GC() {
	_, err := mp.db.deleteExpiredSessions()
	if err != nil {
		panic(err)
	}
}
