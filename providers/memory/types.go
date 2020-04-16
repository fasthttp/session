package memory

import (
	"github.com/fasthttp/session/v2"
)

// Config configuration of provider
type Config struct{}

// Provider backend manager
type Provider struct {
	config Config
	db     *session.Dict
}
