package redis

import (
	"sync"

	"github.com/fasthttp/session"
	"github.com/valyala/fasthttp"
)

var storePool = sync.Pool{
	New: func() interface{} {
		return new(Store)
	},
}

func acquireStore() *Store {
	return storePool.Get().(*Store)
}

func releaseStore(store *Store) {
	store.Reset()
	storePool.Put(store)
}

// NewStore new redis store
func NewStore(sessionID []byte, data *session.Dict) *Store {
	store := acquireStore()
	store.Init(sessionID, data)

	return store
}

// Save save store
func (rs *Store) Save(ctx *fasthttp.RequestCtx) error {
	defer releaseStore(rs)

	data := rs.GetAll()
	b, err := provider.config.SerializeFunc(&data)
	if err != nil {
		return err
	}

	conn := provider.redisPool.Get()
	_, err = conn.Do("SETEX", provider.getRedisSessionKey(rs.GetSessionID()), provider.maxLifeTime, string(b))
	conn.Close()

	return err
}
