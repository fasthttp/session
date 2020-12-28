package memory

import (
	"time"

	"github.com/authelia/session/v2"
)

// Config provider settings
type Config struct{}

// Provider backend manager
type Provider struct {
	config Config
	db     *session.Dict
}

type item struct {
	data           []byte
	lastActiveTime int64
	expiration     time.Duration
}
