package postgres

import (
	"time"

	"github.com/fasthttp/session"
	"github.com/valyala/fasthttp"
)

// session postgres store

// NewPostgresStore new default postgres store
func NewPostgresStore(sessionID string) *Store {
	postgresStore := &Store{}
	postgresStore.Init(sessionID, make(map[string]interface{}))
	return postgresStore
}

// NewPostgresStoreData new postgres store data
func NewPostgresStoreData(sessionID string, data map[string]interface{}) *Store {
	postgresStore := &Store{}
	postgresStore.Init(sessionID, data)
	return postgresStore
}

// Store store struct
type Store struct {
	session.Store
}

// Save save store
func (ps *Store) Save(ctx *fasthttp.RequestCtx) error {

	b, err := provider.config.SerializeFunc(ps.GetAll())
	if err != nil {
		return err
	}
	session, err := provider.sessionDao.getSessionBySessionID(ps.GetSessionID())
	if err != nil || len(session) == 0 {
		return nil
	}
	_, err = provider.sessionDao.updateBySessionID(ps.GetSessionID(), string(b), time.Now().Unix())
	return err
}
