package sqlite3

import (
	"errors"
	"reflect"
	"time"

	"github.com/fasthttp/session"
)

// session sqlite3 provider

//  session Table structure
//
//  DROP TABLE IF EXISTS `session`;
//  CREATE TABLE `session` (
//    `session_id` varchar(64) NOT NULL DEFAULT '',
//    `contents` TEXT NOT NULL,
//    `last_active` int(10) NOT NULL DEFAULT '0',
//    PRIMARY KEY (`session_id`),
//  )
//  create index last_active on session (last_active);
//

// ProviderName sqlite provider name
const ProviderName = "sqlite3"

var (
	provider = NewProvider()
	encrypt  = session.NewEncrypt()
)

// Provider provider struct
type Provider struct {
	config      *Config
	sessionDao  *sessionDao
	maxLifeTime int64
}

// NewProvider new sqlite3 provider
func NewProvider() *Provider {
	return &Provider{
		config:     &Config{},
		sessionDao: &sessionDao{},
	}
}

// Init init provider config
func (sp *Provider) Init(lifeTime int64, sqlite3Config session.ProviderConfig) error {
	if sqlite3Config.Name() != ProviderName {
		return errors.New("session sqlite3 provider init error, config must sqlite3 config")
	}
	vc := reflect.ValueOf(sqlite3Config)
	rc := vc.Interface().(*Config)
	sp.config = rc
	sp.maxLifeTime = lifeTime

	// check config
	if sp.config.DBPath == "" {
		return errors.New("session sqlite3 provider init error, config DBPath not empty")
	}
	// init config serialize func
	if sp.config.SerializeFunc == nil {
		sp.config.SerializeFunc = encrypt.Base64Encode
	}
	if sp.config.UnSerializeFunc == nil {
		sp.config.UnSerializeFunc = encrypt.Base64Decode
	}
	// init sessionDao
	sessionDao, err := newSessionDao(sp.config.DBPath, sp.config.TableName)
	if err != nil {
		return err
	}
	sessionDao.sqlite3Conn.SetMaxOpenConns(sp.config.SetMaxIdleConn)
	sessionDao.sqlite3Conn.SetMaxIdleConns(sp.config.SetMaxIdleConn)

	sp.sessionDao = sessionDao
	return sessionDao.sqlite3Conn.Ping()
}

// NeedGC not need gc
func (sp *Provider) NeedGC() bool {
	return true
}

// GC session sqlite3 provider not need garbage collection
func (sp *Provider) GC() {
	sp.sessionDao.deleteSessionByMaxLifeTime(sp.maxLifeTime)
}

// ReadStore read session store by session id
func (sp *Provider) ReadStore(sessionID string) (session.SessionStore, error) {

	sessionValue, err := sp.sessionDao.getSessionBySessionID(sessionID)
	if err != nil {
		return nil, err
	}
	if len(sessionValue) == 0 {
		_, err := sp.sessionDao.insert(sessionID, "", time.Now().Unix())
		if err != nil {
			return nil, err
		}
		return NewSqLite3Store(sessionID), nil
	}
	if len(sessionValue["contents"]) == 0 {
		return NewSqLite3Store(sessionID), nil
	}

	data, err := sp.config.UnSerializeFunc(sessionValue["contents"])
	if err != nil {
		return nil, err
	}

	return NewSqLite3StoreData(sessionID, data), nil
}

// Regenerate regenerate session
func (sp *Provider) Regenerate(oldSessionId string, sessionID string) (session.SessionStore, error) {

	sessionValue, err := sp.sessionDao.getSessionBySessionID(oldSessionId)
	if err != nil {
		return nil, err
	}
	if len(sessionValue) == 0 {
		// old sessionID not exists, insert new sessionID
		_, err := sp.sessionDao.insert(sessionID, "", time.Now().Unix())
		if err != nil {
			return nil, err
		}
		return NewSqLite3Store(sessionID), nil
	}

	// delete old session
	_, err = sp.sessionDao.deleteBySessionID(oldSessionId)
	if err != nil {
		return nil, err
	}
	// insert new session
	_, err = sp.sessionDao.insert(sessionID, string(sessionValue["contents"]), time.Now().Unix())
	if err != nil {
		return nil, err
	}

	return sp.ReadStore(sessionID)
}

// Destroy destroy session by sessionID
func (sp *Provider) Destroy(sessionID string) error {
	_, err := sp.sessionDao.deleteBySessionID(sessionID)
	return err
}

// Count session values count
func (sp *Provider) Count() int {
	return sp.sessionDao.countSessions()
}

// register session provider
func init() {
	session.Register(ProviderName, provider)
}
