package sqlite3

import (
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
	}
}

// Init init provider config
func (sp *Provider) Init(lifeTime int64, cfg session.ProviderConfig) error {
	if cfg.Name() != ProviderName {
		return errInvalidProviderConfig
	}

	sp.config = cfg.(*Config)
	sp.maxLifeTime = lifeTime

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
	store := NewStore(sessionID)

	row, err := sp.db.getSessionBySessionID(sessionID)

	if row.sessionID != "" { // Exist
		err := sp.config.UnSerializeFunc(gotils.S2B(row.contents), store.GetDataPointer())
		if err != nil {
			return nil, err
		}

	} else { // Not exist
		_, err := sp.db.insert(sessionID, nil, time.Now().Unix())
		if err != nil {
			return nil, err
		}
	}

	releaseDBRow(row)

	return store, err
}

// Regenerate regenerate session
func (sp *Provider) Regenerate(oldID, newID []byte) (session.Storer, error) {
	store := NewStore(newID)

	row, err := sp.db.getSessionBySessionID(oldID)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()

	if row.sessionID != "" { // Exists
		_, err = sp.db.regenerate(oldID, newID, now)
		if err != nil {
			return nil, err
		}

		err = sp.config.UnSerializeFunc(gotils.S2B(row.contents), store.GetDataPointer())
		if err != nil {
			return nil, err
		}

	} else { // Not exist
		_, err = sp.db.insert(newID, nil, now)
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
	_, err := sp.db.deleteSessionByMaxLifeTime(sp.maxLifeTime)
	if err != nil {
		panic(err)
	}
}

// register session provider
func init() {
	session.Register(ProviderName, provider)
}
