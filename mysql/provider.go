package mysql

import (
	"time"

	"github.com/fasthttp/session"
	"github.com/valyala/bytebufferpool"
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

// ReadStore read session store by session id
func (mp *Provider) ReadStore(sessionID []byte) (session.Storer, error) {
	store := NewStore(sessionID)

	row, err := mp.db.getSessionBySessionID(sessionID)

	if row.sessionID != "" { // Exist
		buff := bytebufferpool.Get()
		buff.SetString(row.contents)
		err := mp.config.UnSerializeFunc(buff.Bytes(), store.GetData())
		bytebufferpool.Put(buff)

		if err != nil {
			return nil, err
		}

	} else { // Not exist
		_, err := mp.db.insert(sessionID, nil, time.Now().Unix())
		if err != nil {
			return nil, err
		}
	}

	releaseDBRow(row)

	return store, err
}

// Regenerate regenerate session
func (mp *Provider) Regenerate(oldID, newID []byte) (session.Storer, error) {
	store := NewStore(newID)

	row, err := mp.db.getSessionBySessionID(oldID)
	now := time.Now().Unix()

	if row.sessionID != "" { // Exists
		_, err = mp.db.regenerate(oldID, newID, now)
		if err != nil {
			return nil, err
		}

		buff := bytebufferpool.Get()
		buff.SetString(row.contents)
		err := mp.config.UnSerializeFunc(buff.Bytes(), store.GetData())
		bytebufferpool.Put(buff)

		if err != nil {
			return nil, err
		}

	} else { // Not exist
		_, err := mp.db.insert(newID, nil, now)
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

// GC session mysql provider not need garbage collection
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
