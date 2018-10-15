package session

// Init init store data and sessionID
func (s *Store) Init(sessionID []byte) {
	s.sessionID = sessionID

	if s.data == nil { // Ensure the store always has a valid pointer of Dict
		s.data = new(Dict)
	}
}

// Save save store
func (s *Store) Save() error {
	return nil
}

// Get get data by key
func (s *Store) Get(key string) interface{} {
	return s.data.Get(key)
}

// GetBytes get data by key
func (s *Store) GetBytes(key []byte) interface{} {
	return s.data.GetBytes(key)
}

// GetAll get all data
func (s *Store) GetAll() Dict {
	return *s.data
}

// GetDataPointer get pointer of data
func (s *Store) GetDataPointer() *Dict {
	return s.data
}

// Set set data
func (s *Store) Set(key string, value interface{}) {
	s.data.Set(key, value)
}

// SetBytes set data
func (s *Store) SetBytes(key []byte, value interface{}) {
	s.data.SetBytes(key, value)
}

// Delete delete data by key
func (s *Store) Delete(key string) {
	s.data.Del(key)
}

// DeleteBytes delete data by key
func (s *Store) DeleteBytes(key []byte) {
	s.data.DelBytes(key)
}

// Flush flush all data
func (s *Store) Flush() {
	s.data.Reset()
}

// GetSessionID get session id
func (s *Store) GetSessionID() []byte {
	return s.sessionID
}

// SetSessionID set session id
func (s *Store) SetSessionID(id []byte) {
	s.lock.Lock()
	s.sessionID = id
	s.lock.Unlock()
}

// Reset reset store
func (s *Store) Reset() {
	s.sessionID = nil
	s.data.Reset()
}
