package session

import (
	"sync"
	"time"

	"github.com/savsgio/dictpool"
	"github.com/valyala/fasthttp"
)

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

	// gc life time to execute it
	GCLifetime time.Duration

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
	SessionIDGeneratorFunc func() []byte

	// IsSecureFunc should return whether the communication channel is secure
	// in order to set the secure flag to true according to Secure flag.
	IsSecureFunc func(*fasthttp.RequestCtx) bool

	// value cookie length
	cookieLen uint32
}

// Session session struct
type Session struct {
	provider Provider
	config   *Config
	cookie   *cookie

	storePool *sync.Pool
}

// Dict memory store
type Dict struct {
	dictpool.Dict
}

// Store store
type Store struct {
	sessionID         []byte
	data              *Dict
	defaultExpiration time.Duration
	lock              sync.RWMutex
}

// Cookie cookie struct
type cookie struct{}

// Provider provider interface
type Provider interface {
	Get(store *Store) error
	Save(store *Store) error
	Destroy(id []byte) error
	Regenerate(id []byte, newStore *Store) error // the expiration is also reset to original value
	Count() int
	NeedGC() bool
	GC()
}
