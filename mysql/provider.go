package mysql

import (
	"errors"
	"reflect"
	"time"

	"github.com/fasthttp/session"
)

// session mysql provider

// session Table structure
//
// DROP TABLE IF EXISTS `session`;
// CREATE TABLE `session` (
//    `session_id` varchar(64) NOT NULL DEFAULT '' COMMENT 'Session id',
//    `contents` TEXT NOT NULL COMMENT 'Session data',
//    `last_active` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'Last active time',
//    PRIMARY KEY (`session_id`),
//    KEY `last_active` (`last_active`)
// ) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='session table';
//

// ProviderName mysql provider name
const ProviderName = "mysql"

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

// NewProvider new mysql provider
func NewProvider() *Provider {
	return &Provider{
		config:     &Config{},
		values:     session.NewDefaultCCMap(),
		sessionDao: &sessionDao{},
	}
}

// Init init provider config
func (mp *Provider) Init(lifeTime int64, mysqlConfig session.ProviderConfig) error {
	if mysqlConfig.Name() != ProviderName {
		return errors.New("session mysql provider init error, config must mysql config")
	}
	vc := reflect.ValueOf(mysqlConfig)
	rc := vc.Interface().(*Config)
	mp.config = rc
	mp.maxLifeTime = lifeTime

	// check config
	if mp.config.Host == "" {
		return errors.New("session mysql provider init error, config Host not empty")
	}
	if mp.config.Port == 0 {
		return errors.New("session mysql provider init error, config Port not empty")
	}
	// init config serialize func
	if mp.config.SerializeFunc == nil {
		mp.config.SerializeFunc = encrypt.Base64Encode
	}
	if mp.config.UnSerializeFunc == nil {
		mp.config.UnSerializeFunc = encrypt.Base64Decode
	}
	// init sessionDao
	sessionDao, err := newSessionDao(mp.config.getMysqlDSN(), mp.config.TableName)
	if err != nil {
		return err
	}
	sessionDao.mysqlConn.SetMaxOpenConns(mp.config.SetMaxIdleConn)
	sessionDao.mysqlConn.SetMaxIdleConns(mp.config.SetMaxIdleConn)

	mp.sessionDao = sessionDao
	return sessionDao.mysqlConn.Ping()
}

// NeedGC not need gc
func (mp *Provider) NeedGC() bool {
	return true
}

// GC session mysql provider not need garbage collection
func (mp *Provider) GC() {
	mp.sessionDao.deleteSessionByMaxLifeTime(mp.maxLifeTime)
}

// ReadStore read session store by session id
func (mp *Provider) ReadStore(sessionID string) (session.SessionStore, error) {

	sessionValue, err := mp.sessionDao.getSessionBySessionID(sessionID)
	if err != nil {
		return nil, err
	}
	if len(sessionValue) == 0 {
		_, err := mp.sessionDao.insert(sessionID, "", time.Now().Unix())
		if err != nil {
			return nil, err
		}
		return NewMysqlStore(sessionID), nil
	}
	if len(sessionValue["contents"]) == 0 {
		return NewMysqlStore(sessionID), nil
	}

	data, err := mp.config.UnSerializeFunc(sessionValue["contents"])
	if err != nil {
		return nil, err
	}

	return NewMysqlStoreData(sessionID, data), nil
}

// Regenerate regenerate session
func (mp *Provider) Regenerate(oldSessionId string, sessionID string) (session.SessionStore, error) {

	sessionValue, err := mp.sessionDao.getSessionBySessionID(oldSessionId)
	if err != nil {
		return nil, err
	}
	if len(sessionValue) == 0 {
		// old sessionID not exists, insert new sessionID
		_, err := mp.sessionDao.insert(sessionID, "", time.Now().Unix())
		if err != nil {
			return nil, err
		}
		return NewMysqlStore(sessionID), nil
	}

	// delete old session
	_, err = mp.sessionDao.deleteBySessionID(oldSessionId)
	if err != nil {
		return nil, err
	}
	// insert new session
	_, err = mp.sessionDao.insert(sessionID, string(sessionValue["contents"]), time.Now().Unix())
	if err != nil {
		return nil, err
	}

	return mp.ReadStore(sessionID)
}

// Destroy destroy session by sessionID
func (mp *Provider) Destroy(sessionID string) error {
	_, err := mp.sessionDao.deleteBySessionID(sessionID)
	return err
}

// Count session values count
func (mp *Provider) Count() int {
	return mp.sessionDao.countSessions()
}

// register session provider
func init() {
	session.Register(ProviderName, provider)
}
