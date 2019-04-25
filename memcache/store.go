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

	expiration := mcs.GetExpiration()
	if expiration > math.MaxInt32 {
		return errExpirationIsTooBig
	}
	item.Expiration = int32(expiration)

	err = provider.db.Set(item)

	releaseItem(item)

	return err
}

// Init initialize the store
func (mcs *Store) Init(sessionID []byte, expiration time.Duration) {
	mcs.newExpiration = expiration
	mcs.Store.Init(sessionID, expiration)
}

// SetExpiration set the expiration for the session
func (mcs *Store) SetExpiration(expiration time.Duration) error {
	if expiration > math.MaxInt32 {
		return errExpirationIsTooBig
	}
	mcs.newExpiration = expiration
	return nil
}

// GetExpiration get the expiration for the session
func (mcs *Store) GetExpiration() time.Duration {
	return mcs.newExpiration
}

// HasExpirationChanged return whether the expiration has been updated by the user.
func (mcs *Store) HasExpirationChanged() bool {
	return mcs.newExpiration != mcs.Store.GetExpiration()
}
