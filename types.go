package session

import (
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

// Config configuration of session manager
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
	Expiration time.Duration

	// gc life time to execute it
	GCLifetime time.Duration

	// set whether to pass this bar cookie only through HTTPS
	Secure bool

	// allows you to declare if your cookie should be restricted to a first-party or same-site context.
	// possible values: lax, strict, none
	CookieSameSite fasthttp.CookieSameSite

	// sessionID is in url query
	SessionIDInURLQuery bool

	// sessionName in url query
	SessionNameInURLQuery string

	// sessionID is in http header
	SessionIDInHTTPHeader bool

	// sessionName in http header
	SessionNameInHTTPHeader string

	// SessionIDGeneratorFunc should returns a random session id.
	SessionIDGeneratorFunc func() []byte

	// IsSecureFunc should return whether the communication channel is secure
	// in order to set the secure flag to true according to Secure flag.
	IsSecureFunc func(*fasthttp.RequestCtx) bool

	// EncodeFunc session value serialize func
	EncodeFunc func(src Dict) ([]byte, error)

	// DecodeFunc session value unSerialize func
	DecodeFunc func(dst *Dict, src []byte) error

	// Logger
	Logger Logger

	// value cookie length
	cookieLen uint32
}

// Session manages the users sessions
type Session struct {
	provider Provider
	config   Config
	cookie   *cookie
	log      Logger

	storePool  sync.Pool
	stopGCChan chan struct{}
}

// Store represents the user session
type Store struct {
	sessionID         []byte
	data              Dict
	defaultExpiration time.Duration
	lock              sync.RWMutex
}

type cookie struct{}

type Logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}

// Provider interface implemented by providers
type Provider interface {
	Get(id []byte) ([]byte, error)
	Save(id, data []byte, expiration time.Duration) error
	Destroy(id []byte) error
	Regenerate(id, newID []byte, expiration time.Duration) error
	Count() int
	NeedGC() bool
	GC() error
}
