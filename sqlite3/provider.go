package sqlite3

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

// NewProvider new sqlite3 provider
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

func (sp *Provider) acquireStore(sessionID []byte, expiration time.Duration) *Store {
	store := sp.storePool.Get().(*Store)
	store.Init(sessionID, expiration)

	return store
}

func (sp *Provider) releaseStore(store *Store) {
	store.Reset()
	sp.storePool.Put(store)
}

// Init init provider config
func (sp *Provider) Init(expiration time.Duration, cfg session.ProviderConfig) error {
	if cfg.Name() != ProviderName {
		return errInvalidProviderConfig
	}

	sp.config = cfg.(*Config)
	sp.expiration = expiration

	if sp.config.DBPath == "" {
		return errConfigDBPathEmpty
	}

	if sp.config.SerializeFunc == nil {
		sp.config.SerializeFunc = encrypt.Base64Encode
	}
	if sp.config.UnSerializeFunc == nil {
		sp.config.UnSerializeFunc = encrypt.Base64Decode
	}

	var err error
	sp.db, err = NewDao("sqlite3", sp.config.DBPath, sp.config.TableName)
	if err != nil {
		return err
	}
	sp.db.Connection.SetMaxOpenConns(sp.config.SetMaxIdleConn)
	sp.db.Connection.SetMaxIdleConns(sp.config.SetMaxIdleConn)

	return sp.db.Connection.Ping()
}

// Get read session store by session id
func (sp *Provider) Get(sessionID []byte) (session.Storer, error) {
	store := sp.acquireStore(sessionID, sp.expiration)

	row, err := sp.db.getSessionBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	if row.sessionID != "" { // Exist
		err := sp.config.UnSerializeFunc(store.DataPointer(), gotils.S2B(row.contents))
		if err != nil {
			return nil, err
		}

	} else { // Not exist
		_, err := sp.db.insert(sessionID, nil, time.Now().Unix(), sp.expiration)
		if err != nil {
			return nil, err
		}

	}

	releaseDBRow(row)

	return store, nil
}

// Put put store into the pool.
func (sp *Provider) Put(store session.Storer) {
	sp.releaseStore(store.(*Store))
}

// Regenerate regenerate session
func (sp *Provider) Regenerate(oldID, newID []byte) (session.Storer, error) {
	store := sp.acquireStore(newID, sp.expiration)

	row, err := sp.db.getSessionBySessionID(oldID)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()

	if row.sessionID != "" { // Exists
		_, err = sp.db.regenerate(oldID, newID, now, sp.expiration)
		if err != nil {
			return nil, err
		}

		err = sp.config.UnSerializeFunc(store.DataPointer(), gotils.S2B(row.contents))
		if err != nil {
			return nil, err
		}

	} else { // Not exist
		_, err = sp.db.insert(newID, nil, now, sp.expiration)
		if err != nil {
			return nil, err
		}
	}

	releaseDBRow(row)

	return store, nil
}

// Destroy destroy session by sessionID
func (sp *Provider) Destroy(sessionID []byte) error {
	_, err := sp.db.deleteBySessionID(sessionID)
	return err
}

// Count session values count
func (sp *Provider) Count() int {
	return sp.db.countSessions()
}

// NeedGC need gc
func (sp *Provider) NeedGC() bool {
	return true
}

// GC session garbage collection
func (sp *Provider) GC() {
	_, err := sp.db.deleteExpiredSessions()
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
