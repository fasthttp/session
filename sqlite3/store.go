package sqlite3

import (
	"time"
)

// Save save store
func (ss *Store) Save() error {
	data := ss.GetAll()
	value, err := provider.config.SerializeFunc(data)
	if err != nil {
		return err
	}

	_, err = provider.db.updateBySessionID(ss.GetSessionID(), value, time.Now().Unix())

	return err
}
