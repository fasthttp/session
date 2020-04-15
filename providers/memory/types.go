package memory

import (
	"github.com/fasthttp/session"
)

// Config session memory configuration
type Config struct{}

// Provider provider struct
type Provider struct {
	config Config
	db     *session.Dict
}
