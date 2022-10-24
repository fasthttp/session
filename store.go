package session

import (
	"time"

	"github.com/savsgio/gotils/strconv"
)

func newDictValue() Dict {
	return Dict{
		KV: make(map[string]interface{}),
	}
}

// NewStore returns a new empty store
func NewStore() *Store {
	return &Store{
		data: newDictValue(),
	}
}

// Get returns a value from the given key
func (s *Store) Get(key string) interface{} {
	return s.data.KV[key]
}

// GetBytes returns a value from the given key
func (s *Store) GetBytes(key []byte) interface{} {
	return s.Get(strconv.B2S(key))
}

// GetAll returns all stored values
func (s *Store) GetAll() Dict {
	return s.data
}

// Ptr returns the internal store pointer
func (s *Store) Ptr() *Dict {
	return &s.data
}

// Set saves a value for the given key
func (s *Store) Set(key string, value interface{}) {
	s.data.KV[key] = value
}

// SetBytes saves a value for the given key
func (s *Store) SetBytes(key []byte, value interface{}) {
	s.Set(strconv.B2S(key), value)
}

// Delete deletes a value from the given key
func (s *Store) Delete(key string) {
	delete(s.data.KV, key)
}

// DeleteBytes deletes a value from the given key
func (s *Store) DeleteBytes(key []byte) {
	s.Delete(strconv.B2S(key))
}

// Flush removes all stored values
func (s *Store) Flush() {
	for k := range s.data.KV {
		delete(s.data.KV, k)
	}
}

// GetSessionID returns the session id
func (s *Store) GetSessionID() []byte {
	s.lock.RLock()
	id := s.sessionID
	s.lock.RUnlock()

	return id
}

// SetSessionID sets the session id
func (s *Store) SetSessionID(id []byte) {
	s.lock.Lock()
	s.sessionID = id
	s.lock.Unlock()
}

// HasExpirationChanged checks wether the expiration has been changed
func (s *Store) HasExpirationChanged() bool {
	return s.Get(expirationAttrKey) != nil
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
	s.Flush()
	s.sessionID = s.sessionID[:0]
	s.defaultExpiration = 0
}
