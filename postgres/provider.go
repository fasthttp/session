package postgres

import (
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
	}
}

// Init init provider config
func (pp *Provider) Init(lifeTime int64, cfg session.ProviderConfig) error {
	if cfg.Name() != ProviderName {
		return errInvalidProviderConfig
	}

	pp.config = cfg.(*Config)
	pp.maxLifeTime = lifeTime

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

// ReadStore read session store by session id
func (pp *Provider) ReadStore(sessionID []byte) (session.Storer, error) {
	store := NewStore(sessionID)

	row, err := pp.db.getSessionBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	if row.sessionID != "" { // Exist
		err = pp.config.UnSerializeFunc(gotils.S2B(row.contents), store.GetDataPointer())
		if err != nil {
			return nil, err
		}

	} else { // Not exist
		_, err = pp.db.insert(sessionID, nil, time.Now().Unix())
		if err != nil {
			return nil, err
		}
	}

	releaseDBRow(row)

	return store, nil
}

// Regenerate regenerate session
func (pp *Provider) Regenerate(oldID, newID []byte) (session.Storer, error) {
	store := NewStore(newID)

	row, err := pp.db.getSessionBySessionID(oldID)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()

	if row.sessionID != "" { // Exists
		_, err = pp.db.regenerate(oldID, newID, now)
		if err != nil {
			return nil, err
		}

		err = pp.config.UnSerializeFunc(gotils.S2B(row.contents), store.GetDataPointer())
		if err != nil {
			return nil, err
		}

	} else { // Not exist
		_, err = pp.db.insert(newID, nil, now)
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
	_, err := pp.db.deleteSessionByMaxLifeTime(pp.maxLifeTime)
	if err != nil {
		panic(err)
	}
}

// register session provider
func init() {
	session.Register(ProviderName, provider)
}
