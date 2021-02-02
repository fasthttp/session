package session

import (
	"github.com/savsgio/gotils/bytes"
	"github.com/valyala/fasthttp"
)

// NewDefaultConfig returns a new default configuration
func NewDefaultConfig() Config {
	config := Config{
		CookieName:              defaultSessionKeyName,
		Domain:                  defaultDomain,
		Expiration:              defaultExpiration,
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

	// default isSecureFunc
	config.IsSecureFunc = config.defaultIsSecureFunc

	return config
}

func (c *Config) defaultSessionIDGenerator() []byte {
	return bytes.Rand(make([]byte, c.cookieLen))
}

func (c *Config) defaultIsSecureFunc(ctx *fasthttp.RequestCtx) bool {
	return ctx.IsTLS()
}
