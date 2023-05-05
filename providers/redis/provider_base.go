package redis

import (
	"context"
	"time"

	"github.com/valyala/bytebufferpool"
)

func (p *Provider) getRedisSessionKey(sessionID []byte) string {
	key := bytebufferpool.Get()
	key.SetString(p.keyPrefix)
	key.WriteString(":")
	key.Write(sessionID)

	keyStr := key.String()

	bytebufferpool.Put(key)

	return keyStr
}

// Save saves the session data and expiration from the given session id
func (p *Provider) Save(id, data []byte, expiration time.Duration) error {
	key := p.getRedisSessionKey(id)

	return p.db.Set(context.Background(), key, data, expiration).Err()
}

// Regenerate updates the session id and expiration with the new session id
// of the given current session id
func (p *Provider) Regenerate(id, newID []byte, expiration time.Duration) error {
	key := p.getRedisSessionKey(id)
	newKey := p.getRedisSessionKey(newID)

	exists, err := p.db.Exists(context.Background(), key).Result()
	if err != nil {
		return err
	}

	if exists > 0 { // Exist
		if err = p.db.Rename(context.Background(), key, newKey).Err(); err != nil {
			return err
		}

		if err = p.db.Expire(context.Background(), newKey, expiration).Err(); err != nil {
			return err
		}
	}

	return nil
}

// Destroy destroys the session from the given id
func (p *Provider) Destroy(id []byte) error {
	key := p.getRedisSessionKey(id)

	return p.db.Del(context.Background(), key).Err()
}

// Count returns the total of stored sessions
func (p *Provider) Count() int {
	reply, err := p.db.Keys(context.Background(), p.getRedisSessionKey(all)).Result()
	if err != nil {
		return 0
	}

	return len(reply)
}

// NeedGC indicates if the GC needs to be run
func (p *Provider) NeedGC() bool {
	return false
}

// GC destroys the expired sessions
func (p *Provider) GC() error {
	return nil
}
