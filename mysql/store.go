package mysql

import (
	"time"

	"github.com/fasthttp/session"
	"github.com/valyala/fasthttp"
)

// session mysql store

// NewMysqlStore new default mysql store
func NewMysqlStore(sessionID string) *Store {
	mysqlStore := &Store{}
	mysqlStore.Init(sessionID, make(map[string]interface{}))
	return mysqlStore
}

// NewMysqlStoreData new mysql store data
func NewMysqlStoreData(sessionID string, data map[string]interface{}) *Store {
	mysqlStore := &Store{}
	mysqlStore.Init(sessionID, data)
	return mysqlStore
}

// Store store struct
type Store struct {
	session.Store
}

// Save save store
func (ms *Store) Save(ctx *fasthttp.RequestCtx) error {

	b, err := provider.config.SerializeFunc(ms.GetAll())
	if err != nil {
		return err
	}
	session, err := provider.sessionDao.getSessionBySessionID(ms.GetSessionID())
	if err != nil || len(session) == 0 {
		return nil
	}
	_, err = provider.sessionDao.updateBySessionID(ms.GetSessionID(), string(b), time.Now().Unix())
	return err
}
