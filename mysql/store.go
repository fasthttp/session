package mysql

import (
	"time"
)

// Save save store
func (ms *Store) Save() error {
	data := ms.GetAll()
	value, err := provider.config.SerializeFunc(data)
	if err != nil {
		return err
	}

	_, err = provider.db.updateBySessionID(ms.GetSessionID(), value, time.Now().Unix())

	return err
}
