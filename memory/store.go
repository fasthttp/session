package memory

import (
	"time"
)

// Save save store
func (ms *Store) Save() error {
	ms.lock.Lock()
	ms.lastActiveTime = time.Now().Unix()
	ms.lock.Unlock()

	return nil
}
