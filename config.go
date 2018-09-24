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
		SessionLifetime:         60,
		Secure:                  true,
		SessionIDInURLQuery:     false,
		SessionNameInURLQuery:   defaultSessionKeyName,
		SessionIDInHTTPHeader:   false,
		SessionNameInHTTPHeader: defaultSessionKeyName,
	}

	// default sessionIdGeneratorFunc
	config.SessionIDGeneratorFunc = config.defaultSessionIDGenerator

	return config
}

// SessionIDGenerator sessionID generator
func (c *Config) SessionIDGenerator() []byte {
	return c.SessionIDGeneratorFunc()
}

// default sessionID generator => uuid
func (c *Config) defaultSessionIDGenerator() []byte {
	b := make([]byte, c.cookieLen)

	for i := 0; i < int(c.cookieLen); i++ {
		pos := fastrand.Uint32n(c.cookieLen)
		b[i] = cookieCharset[pos]
	}

	return b
}

// Encode encode cookie value
func (c *Config) Encode(cookieValue []byte) []byte {
	encode := c.EncodeFunc
	if encode != nil {
		newValue, err := encode(cookieValue)
		if err != nil {
			newValue = nil
		}

		return newValue
	}

	return cookieValue
}

// Decode decode cookie value
func (c *Config) Decode(cookieValue []byte) []byte {
	if len(cookieValue) == 0 {
		return nil
	}

	decode := c.DecodeFunc
	if decode != nil {
		newValue, err := decode(cookieValue)
		if err != nil {
			newValue = nil
		}

		return newValue
	}

	return cookieValue
}
