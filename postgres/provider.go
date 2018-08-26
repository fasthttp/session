package postgres

import (
	"errors"
	"reflect"
	"time"

	"github.com/fasthttp/session"
)

// session postgres provider

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

// ProviderName postgres provider name
const ProviderName = "postgres"

var (
	provider = NewProvider()
	encrypt  = session.NewEncrypt()
)

// Provider provider struct
type Provider struct {
	config      *Config
	values      *session.CCMap
	sessionDao  *sessionDao
	maxLifeTime int64
}

// NewProvider new postgres provider
func NewProvider() *Provider {
	return &Provider{
		config:     &Config{},
		values:     session.NewDefaultCCMap(),
		sessionDao: &sessionDao{},
	}
}

// Init init provider config
func (pp *Provider) Init(lifeTime int64, postgresConfig session.ProviderConfig) error {
	if postgresConfig.Name() != ProviderName {
		return errors.New("session postgres provider init error, config must postgres config")
	}
	vc := reflect.ValueOf(postgresConfig)
	rc := vc.Interface().(*Config)
	pp.config = rc
	pp.maxLifeTime = lifeTime

	// check config
	if pp.config.Host == "" {
		return errors.New("session postgres provider init error, config Host not empty")
	}
	if pp.config.Port == 0 {
		return errors.New("session postgres provider init error, config Port not empty")
	}
	// init config serialize func
	if pp.config.SerializeFunc == nil {
		pp.config.SerializeFunc = encrypt.Base64Encode
	}
	if pp.config.UnSerializeFunc == nil {
		pp.config.UnSerializeFunc = encrypt.Base64Decode
	}
	// init sessionDao
	sessionDao, err := newSessionDao(pp.config.Database, pp.config.TableName)
	if err != nil {
		return err
	}
	sessionDao.postgresConn.SetMaxOpenConns(pp.config.SetMaxIdleConn)
	sessionDao.postgresConn.SetMaxIdleConns(pp.config.SetMaxIdleConn)

	pp.sessionDao = sessionDao
	return sessionDao.postgresConn.Ping()
}

// NeedGC not need gc
func (pp *Provider) NeedGC() bool {
	return true
}

// GC session postgres provider not need garbage collection
func (pp *Provider) GC() {
	pp.sessionDao.deleteSessionByMaxLifeTime(pp.maxLifeTime)
}

// ReadStore read session store by session id
func (pp *Provider) ReadStore(sessionID string) (session.SessionStore, error) {

	sessionValue, err := pp.sessionDao.getSessionBySessionID(sessionID)
	if err != nil {
		return nil, err
	}
	if len(sessionValue) == 0 {
		_, err := pp.sessionDao.insert(sessionID, "", time.Now().Unix())
		if err != nil {
			return nil, err
		}
		return NewPostgresStore(sessionID), nil
	}
	if len(sessionValue["contents"]) == 0 {
		return NewPostgresStore(sessionID), nil
	}

	data, err := pp.config.UnSerializeFunc(sessionValue["contents"])
	if err != nil {
		return nil, err
	}

	return NewPostgresStoreData(sessionID, data), nil
}

// Regenerate regenerate session
func (pp *Provider) Regenerate(oldSessionID string, sessionID string) (session.SessionStore, error) {

	sessionValue, err := pp.sessionDao.getSessionBySessionID(oldSessionID)
	if err != nil {
		return nil, err
	}
	if len(sessionValue) == 0 {
		// old sessionID not exists, insert new sessionID
		_, err := pp.sessionDao.insert(sessionID, "", time.Now().Unix())
		if err != nil {
			return nil, err
		}
		return NewPostgresStore(sessionID), nil
	}

	// delete old session
	_, err = pp.sessionDao.deleteBySessionID(oldSessionID)
	if err != nil {
		return nil, err
	}
	// insert new session
	_, err = pp.sessionDao.insert(sessionID, string(sessionValue["contents"]), time.Now().Unix())
	if err != nil {
		return nil, err
	}

	return pp.ReadStore(sessionID)
}

// Destroy destroy session by sessionID
func (pp *Provider) Destroy(sessionID string) error {
	_, err := pp.sessionDao.deleteBySessionID(sessionID)
	return err
}

// Count session values count
func (pp *Provider) Count() int {
	return pp.sessionDao.countSessions()
}

// register session provider
func init() {
	session.Register(ProviderName, provider)
}
