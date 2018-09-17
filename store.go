package session

import (
	"github.com/savsgio/dictpool"
	"github.com/valyala/fasthttp"
)

// SessionStore session store struct
type SessionStore interface {
	Save(*fasthttp.RequestCtx) error
	Get(key string) interface{}
	GetAll() *dictpool.Dict
	Set(key string, value interface{})
	Delete(key string)
	Flush()
	GetSessionID() string
}

// Store store
type Store struct {
	sessionID string
	data      dictpool.Dict
}

// Init init store data and sessionID
func (s *Store) Init(sessionID string, data *dictpool.Dict) {
	s.sessionID = sessionID
	s.data = dictpool.Dict{}

	if data != nil {
		for _, kv := range data.D {
			s.data.SetBytes(kv.Key, kv.Value)
		}
	}
}

// Get get data by key
func (s *Store) Get(key string) interface{} {
	return s.data.Get(key)
}

// GetAll get all data
func (s *Store) GetAll() *dictpool.Dict {
	return &s.data
}

// Set set data
func (s *Store) Set(key string, value interface{}) {
	s.data.Set(key, value)
}

// Delete delete data by key
func (s *Store) Delete(key string) {
	s.data.Del(key)
}

// Flush flush all data
func (s *Store) Flush() {
	s.data.Reset()
}

// GetSessionID get session id
func (s *Store) GetSessionID() string {
	return s.sessionID
}
