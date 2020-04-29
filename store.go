package session

import (
	"fmt"
	"time"

	"github.com/savsgio/gotils"
)

var expirationAttrKey = fmt.Sprintf("__store:expiration:%s__", gotils.RandBytes(make([]byte, 5)))

// NewStore returns a new empty store
func NewStore() *Store {
	return &Store{
		data: new(Dict),
	}
}

// Get returns a value from the given key
func (s *Store) Get(key string) interface{} {
	return s.data.Get(key)
}

// GetBytes returns a value from the given key
func (s *Store) GetBytes(key []byte) interface{} {
	return s.data.GetBytes(key)
}

// GetAll returns all stored values
func (s *Store) GetAll() Dict {
	return *s.data
}

// Ptr returns the internal store pointer
func (s *Store) Ptr() *Dict {
	return s.data
}

// Set saves a value for the given key
func (s *Store) Set(key string, value interface{}) {
	s.data.Set(key, value)
}

// SetBytes saves a value for the given key
func (s *Store) SetBytes(key []byte, value interface{}) {
	s.data.SetBytes(key, value)
}

// Delete deletes a value from the given key
func (s *Store) Delete(key string) {
	s.data.Del(key)
}

// DeleteBytes deletes a value from the given key
func (s *Store) DeleteBytes(key []byte) {
	s.data.DelBytes(key)
}

// Flush removes all stored values
func (s *Store) Flush() {
	s.data.Reset()
}

// GetSessionID returns the session id
func (s *Store) GetSessionID() []byte {
	return s.sessionID
}

// SetSessionID sets the session id
func (s *Store) SetSessionID(id []byte) {
	s.lock.Lock()
	s.sessionID = id
	s.lock.Unlock()
}

// HasExpirationChanged checks wether the expiration has been changed
func (s *Store) HasExpirationChanged() bool {
	return s.data.Has(expirationAttrKey)
}

// GetExpiration returns the expiration for current session
func (s *Store) GetExpiration() time.Duration {
	expiration, ok := s.Get(expirationAttrKey).(int64)
	if !ok {
		return s.defaultExpiration
	}

	return time.Duration(expiration)
}

// SetExpiration sets the expiration for current session
func (s *Store) SetExpiration(expiration time.Duration) error {
	s.Set(expirationAttrKey, int64(expiration))

	return nil
}

// Reset resets the store
func (s *Store) Reset() {
	s.data.Reset()
	s.sessionID = s.sessionID[:0]
	s.defaultExpiration = 0
}
