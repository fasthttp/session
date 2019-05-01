package memcache

import (
	"math"
	"time"
)

// Save save store
func (mcs *Store) Save() error {
	data := mcs.GetAll()
	value, err := provider.config.SerializeFunc(data)
	if err != nil {
		return err
	}

	item := acquireItem()
	item.Key = provider.getMemCacheSessionKey(mcs.GetSessionID())
	item.Value = value
	item.Expiration = int32(mcs.GetExpiration() / time.Second)

	err = provider.db.Set(item)

	releaseItem(item)

	return err
}

// SetExpiration set the expiration for the session
func (mcs *Store) SetExpiration(expiration time.Duration) error {
	if expiration/time.Second > math.MaxInt32 {
		return errExpirationIsTooBig
	}
	return mcs.Store.SetExpiration(expiration)
}
