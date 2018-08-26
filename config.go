package session

import (
	"time"

	"github.com/satori/go.uuid"
)

var (
	defaultCookieName = "sessionid"
	defaultExpires    = time.Hour * 2
	defaultGCLifetime = int64(3)
)

// NewDefaultConfig return new default config
func NewDefaultConfig() *Config {
	config := &Config{
		CookieName:              defaultCookieName,
		Domain:                  "",
		Expires:                 defaultExpires,
		GCLifetime:              defaultGCLifetime,
		SessionLifetime:         60,
		Secure:                  true,
		SessionIDInURLQuery:     false,
		SessionNameInURLQuery:   "",
		SessionIDInHTTPHeader:   false,
		SessionNameInHTTPHeader: "",
	}

	// default sessionIdGeneratorFunc
	config.SessionIDGeneratorFunc = config.defaultSessionIDGenerator

	return config
}

// Config config struct
type Config struct {

	// cookie name
	CookieName string

	// cookie domain
	Domain string

	// If you want to delete the cookie when the browser closes, set it to -1.
	//
	//  0 means no expire, (24 years)
	// -1 means when browser closes
	// >0 is the time.Duration which the session cookies should expire.
	Expires time.Duration

	// gc life time(s)
	GCLifetime int64

	// session life time(s)
	SessionLifetime int64

	// set whether to pass this bar cookie only through HTTPS
	Secure bool

	// sessionID is in url query
	SessionIDInURLQuery bool

	// sessionName in url query
	SessionNameInURLQuery string

	// sessionID is in http header
	SessionIDInHTTPHeader bool

	// sessionName in http header
	SessionNameInHTTPHeader string

	// SessionIDGeneratorFunc should returns a random session id.
	SessionIDGeneratorFunc func() string

	// Encode the cookie value if not nil.
	EncodeFunc func(cookieValue string) (string, error)

	// Decode the cookie value if not nil.
	DecodeFunc func(cookieValue string) (string, error)
}

// SessionIDGenerator sessionID generator
func (c *Config) SessionIDGenerator() string {
	sessionIDGenerator := c.SessionIDGeneratorFunc
	if sessionIDGenerator == nil {
		return c.defaultSessionIDGenerator()
	}

	return sessionIDGenerator()
}

// default sessionID generator => uuid
func (c *Config) defaultSessionIDGenerator() string {
	return uuid.NewV4().String()
}

// Encode encode cookie value
func (c *Config) Encode(cookieValue string) string {
	encode := c.EncodeFunc
	if encode != nil {
		newVal, err := encode(cookieValue)
		if err == nil {
			cookieValue = newVal
		} else {
			cookieValue = ""
		}
	}
	return cookieValue
}

// Decode decode cookie value
func (c *Config) Decode(cookieValue string) string {
	if cookieValue == "" {
		return ""
	}
	decode := c.DecodeFunc
	if decode != nil {
		newVal, err := decode(cookieValue)
		if err == nil {
			cookieValue = newVal
		} else {
			cookieValue = ""
		}
	}
	return cookieValue
}
