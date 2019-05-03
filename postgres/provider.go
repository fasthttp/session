package postgres

import (
	"sync"
	"time"

	"github.com/fasthttp/session"
	"github.com/savsgio/gotils"
)

var (
	provider = NewProvider()
	encrypt  = session.NewEncrypt()
)

// NewProvider new postgres provider
func NewProvider() *Provider {
	return &Provider{
		config: new(Config),
		db:     new(Dao),

		storePool: sync.Pool{
			New: func() interface{} {
				return new(Store)
			},
		},
	}
}

func (pp *Provider) acquireStore(sessionID []byte, expiration time.Duration) *Store {
	store := pp.storePool.Get().(*Store)
	store.Init(sessionID, expiration)

	return store
}

func (pp *Provider) releaseStore(store *Store) {
	store.Reset()
	pp.storePool.Put(store)
}

// Init init provider config
func (pp *Provider) Init(expiration time.Duration, cfg session.ProviderConfig) error {
	if cfg.Name() != ProviderName {
		return errInvalidProviderConfig
	}

	pp.config = cfg.(*Config)
	pp.expiration = expiration

	if pp.config.Host == "" {
		return errConfigHostEmpty
	}
	if pp.config.Port == 0 {
		return errConfigPortZero
	}

	if pp.config.SerializeFunc == nil {
		pp.config.SerializeFunc = encrypt.Base64Encode
	}
	if pp.config.UnSerializeFunc == nil {
		pp.config.UnSerializeFunc = encrypt.Base64Decode
	}

	var err error
	pp.db, err = NewDao("postgres", pp.config.getPostgresDSN(), pp.config.TableName)
	if err != nil {
		return err
	}
	pp.db.Connection.SetMaxOpenConns(pp.config.SetMaxIdleConn)
	pp.db.Connection.SetMaxIdleConns(pp.config.SetMaxIdleConn)

	return pp.db.Connection.Ping()
}

// Get read session store by session id
func (pp *Provider) Get(sessionID []byte) (session.Storer, error) {
	store := pp.acquireStore(sessionID, pp.expiration)

	row, err := pp.db.getSessionBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	if row.sessionID != "" { // Exist
		err = pp.config.UnSerializeFunc(store.DataPointer(), gotils.S2B(row.contents))
		if err != nil {
			return nil, err
		}

	} else { // Not exist
		_, err = pp.db.insert(sessionID, nil, time.Now().Unix(), pp.expiration)
		if err != nil {
			return nil, err
		}

	}

	releaseDBRow(row)

	return store, nil
}

// Put put store into the pool.
func (pp *Provider) Put(store session.Storer) {
	pp.releaseStore(store.(*Store))
}

// Regenerate regenerate session
func (pp *Provider) Regenerate(oldID, newID []byte) (session.Storer, error) {
	store := pp.acquireStore(newID, pp.expiration)

	row, err := pp.db.getSessionBySessionID(oldID)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()

	if row.sessionID != "" { // Exists
		_, err = pp.db.regenerate(oldID, newID, now, pp.expiration)
		if err != nil {
			return nil, err
		}

		err = pp.config.UnSerializeFunc(store.DataPointer(), gotils.S2B(row.contents))
		if err != nil {
			return nil, err
		}

	} else { // Not exist
		_, err = pp.db.insert(newID, nil, now, pp.expiration)
		if err != nil {
			return nil, err
		}
	}

	releaseDBRow(row)

	return store, nil
}

// Destroy destroy session by sessionID
func (pp *Provider) Destroy(sessionID []byte) error {
	_, err := pp.db.deleteBySessionID(sessionID)
	return err
}

// Count session values count
func (pp *Provider) Count() int {
	return pp.db.countSessions()
}

// NeedGC need gc
func (pp *Provider) NeedGC() bool {
	return true
}

// GC session garbage collection
func (pp *Provider) GC() {
	_, err := pp.db.deleteExpiredSessions()
	if err != nil {
		panic(err)
	}
}

// register session provider
func init() {
	err := session.Register(ProviderName, provider)
	if err != nil {
		panic(err)
	}
}
