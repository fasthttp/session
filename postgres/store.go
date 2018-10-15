package postgres

import (
	"time"
)

// Save save store
func (ps *Store) Save() error {
	data := ps.GetAll()
	value, err := provider.config.SerializeFunc(data)
	if err != nil {
		return err
	}

	_, err = provider.db.updateBySessionID(ps.GetSessionID(), value, time.Now().Unix())

	return err
}
