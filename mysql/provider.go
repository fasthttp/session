package mysql

import (
	"time"

	"github.com/fasthttp/session"
	"github.com/savsgio/gotils"
)

var (
	provider = NewProvider()
	encrypt  = session.NewEncrypt()
)

// NewProvider new mysql provider
func NewProvider() *Provider {
	return &Provider{
		config: new(Config),
		db:     new(Dao),
	}
}

// Init init provider config
func (mp *Provider) Init(lifeTime int64, cfg session.ProviderConfig) error {
	if cfg.Name() != ProviderName {
		return errInvalidProviderConfig
	}

	mp.config = cfg.(*Config)
	mp.maxLifeTime = lifeTime

	if mp.config.Host == "" {
		return errConfigHostEmpty
	}
	if mp.config.Port == 0 {
		return errConfigPortZero
	}

	if mp.config.SerializeFunc == nil {
		mp.config.SerializeFunc = encrypt.Base64Encode
	}
	if mp.config.UnSerializeFunc == nil {
		mp.config.UnSerializeFunc = encrypt.Base64Decode
	}

	var err error
	mp.db, err = NewDao("mysql", mp.config.getMysqlDSN(), mp.config.TableName)
	if err != nil {
		return err
	}
	mp.db.Connection.SetMaxOpenConns(mp.config.SetMaxIdleConn)
	mp.db.Connection.SetMaxIdleConns(mp.config.SetMaxIdleConn)

	return mp.db.Connection.Ping()
}

// Get read session store by session id
func (mp *Provider) Get(sessionID []byte) (session.Storer, error) {
	store := NewStore(sessionID)

	row, err := mp.db.getSessionBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	if row.sessionID != "" { // Exist
		err = mp.config.UnSerializeFunc(gotils.S2B(row.contents), store.GetDataPointer())
		if err != nil {
			return nil, err
		}

	} else { // Not exist
		_, err = mp.db.insert(sessionID, nil, time.Now().Unix())
		if err != nil {
			return nil, err
		}
	}

	releaseDBRow(row)

	return store, nil
}

// Regenerate regenerate session
func (mp *Provider) Regenerate(oldID, newID []byte) (session.Storer, error) {
	store := NewStore(newID)

	row, err := mp.db.getSessionBySessionID(oldID)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()

	if row.sessionID != "" { // Exists
		_, err = mp.db.regenerate(oldID, newID, now)
		if err != nil {
			return nil, err
		}

		err = mp.config.UnSerializeFunc(gotils.S2B(row.contents), store.GetDataPointer())
		if err != nil {
			return nil, err
		}

	} else { // Not exist
		_, err = mp.db.insert(newID, nil, now)
		if err != nil {
			return nil, err
		}
	}

	releaseDBRow(row)

	return store, nil
}

// Destroy destroy session by sessionID
func (mp *Provider) Destroy(sessionID []byte) error {
	_, err := mp.db.deleteBySessionID(sessionID)
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
	_, err := mp.db.deleteSessionByMaxLifeTime(mp.maxLifeTime)
	if err != nil {
		panic(err)
	}
}

// register session provider
func init() {
	session.Register(ProviderName, provider)
}
