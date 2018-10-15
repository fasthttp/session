package redis

import (
	"github.com/savsgio/gotils"
)

// Save save store
func (rs *Store) Save() error {
	data := rs.GetAll()
	b, err := provider.config.SerializeFunc(data)
	if err != nil {
		return err
	}

	conn := provider.redisPool.Get()
	_, err = conn.Do("SETEX", provider.getRedisSessionKey(rs.GetSessionID()), provider.maxLifeTime, gotils.B2S(b))
	conn.Close()

	return err
}
