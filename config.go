package session

import (
	"github.com/valyala/fastrand"
)

// NewDefaultConfig return new default configuration
func NewDefaultConfig() *Config {
	config := &Config{
		CookieName:              defaultSessionKeyName,
		Domain:                  "",
		Expires:                 defaultExpires,
		GCLifetime:              defaultGCLifetime,
		SessionLifetime:         defaultSessionLifetime,
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

// default sessionID generator. Returns a random session id
func (c *Config) defaultSessionIDGenerator() []byte {
	b := make([]byte, c.cookieLen)

	for i := 0; i < int(c.cookieLen); i++ {
		pos := fastrand.Uint32n(c.cookieLen)
		b[i] = cookieCharset[pos]
	}

	return b
}
