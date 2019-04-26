package session

import (
	"github.com/savsgio/gotils"
)

// NewDefaultConfig return new default configuration
func NewDefaultConfig() *Config {
	config := &Config{
		CookieName:              defaultSessionKeyName,
		Domain:                  defaultDomain,
		Expires:                 defaultExpires,
		GCLifetime:              defaultGCLifetime,
		Secure:                  defaultSecure,
		SessionIDInURLQuery:     defaultSessionIDInURLQuery,
		SessionNameInURLQuery:   defaultSessionKeyName,
		SessionIDInHTTPHeader:   defaultSessionIDInHTTPHeader,
		SessionNameInHTTPHeader: defaultSessionKeyName,
		cookieLen:               defaultCookieLen,
	}

	// default sessionIdGeneratorFunc
	config.SessionIDGeneratorFunc = config.defaultSessionIDGenerator

	return config
}

func (c *Config) defaultSessionIDGenerator() []byte {
	b := make([]byte, c.cookieLen)

	gotils.RandBytes(b)

	return b
}
