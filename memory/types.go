package memory

import (
	"sync"
	"time"

	"github.com/fasthttp/session"
)

// Config session memory configuration
type Config struct{}

// Provider provider struct
type Provider struct {
	config     *Config
	memoryDB   *session.Dict
	expiration time.Duration

	storePool sync.Pool

	lock sync.RWMutex
}

// Store memory store
type Store struct {
	session.Store

	lastActiveTime int64
	lock           sync.RWMutex
}
