package memory

import (
	"sync"
	"time"
)

// Config provider settings
type Config struct{}

// Provider backend manager
type Provider struct {
	config Config
	db     sync.Map
}

type item struct {
	data           []byte
	lastActiveTime int64
	expiration     time.Duration
}
