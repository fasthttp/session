package sqlite3

import (
	"sync"
	"time"
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

// NewStore new sqlite3 store
func NewStore(sessionID []byte) *Store {
	store := acquireStore()
	store.Init(sessionID)

	return store
}

// Save save store
func (ss *Store) Save() error {
	defer releaseStore(ss)

	data := ss.GetAll()
	value, err := provider.config.SerializeFunc(data)
	if err != nil {
		return err
	}

	_, err = provider.db.updateBySessionID(ss.GetSessionID(), value, time.Now().Unix())

	return err
}
