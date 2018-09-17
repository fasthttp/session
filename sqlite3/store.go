package sqlite3

import (
	"time"

	"github.com/fasthttp/session"
	"github.com/savsgio/dictpool"
	"github.com/valyala/fasthttp"
)

// session sqlite3 store

// NewSqLite3Store new default sqlite3 store
func NewSqLite3Store(sessionID string) *Store {
	sqlite3Store := &Store{}
	sqlite3Store.Init(sessionID, nil)
	return sqlite3Store
}

// NewSqLite3StoreData new sqlite3 store data
func NewSqLite3StoreData(sessionID string, data *dictpool.Dict) *Store {
	sqlite3Store := &Store{}
	sqlite3Store.Init(sessionID, data)
	return sqlite3Store
}

// Store store struct
type Store struct {
	session.Store
}

// Save save store
func (ss *Store) Save(ctx *fasthttp.RequestCtx) error {

	b, err := provider.config.SerializeFunc(ss.GetAll())
	if err != nil {
		return err
	}
	session, err := provider.sessionDao.getSessionBySessionID(ss.GetSessionID())
	if err != nil || len(session) == 0 {
		return nil
	}
	_, err = provider.sessionDao.updateBySessionID(ss.GetSessionID(), string(b), time.Now().Unix())
	return err
}
