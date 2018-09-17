package redis

import (
	"github.com/fasthttp/session"
	"github.com/savsgio/dictpool"
	"github.com/valyala/fasthttp"
)

// session redis store

// NewRedisStore new default redis store
func NewRedisStore(sessionID string) *Store {
	redisStore := &Store{}
	redisStore.Init(sessionID, nil)
	return redisStore
}

// NewRedisStoreData new redis store data
func NewRedisStoreData(sessionID string, data *dictpool.Dict) *Store {
	redisStore := &Store{}
	redisStore.Init(sessionID, data)
	return redisStore
}

// Store store struct
type Store struct {
	session.Store
}

// Save save store
func (rs *Store) Save(ctx *fasthttp.RequestCtx) error {

	b, err := provider.config.SerializeFunc(rs.GetAll())
	if err != nil {
		return err
	}
	conn := provider.redisPool.Get()
	defer conn.Close()
	conn.Do("SETEX", provider.getRedisSessionKey(rs.GetSessionID()), provider.maxLifeTime, string(b))

	return nil
}
