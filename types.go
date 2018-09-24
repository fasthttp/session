package session

import (
	"database/sql"
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
	SessionIDGeneratorFunc func() []byte

	// Encode the cookie value if not nil.
	EncodeFunc func(cookieValue []byte) ([]byte, error)

	// Decode the cookie value if not nil.
	DecodeFunc func(cookieValue []byte) ([]byte, error)

	cookieLen uint32
}

// Dict memory store
type Dict struct {
	dictpool.Dict
}

// Session session struct
type Session struct {
	provider Provider
	config   *Config
	cookie   *Cookie
}

// Dao database connection
type Dao struct {
	Driver     string
	Dsn        string
	Connection *sql.DB
}

// Store store
type Store struct {
	sessionID []byte
	data      Dict
	lock      sync.RWMutex
}

// Encrypt encrypt struct
type Encrypt struct{}

// Cookie cookie struct
type Cookie struct{}

// Storer session store interface
type Storer interface {
	Save(ctx *fasthttp.RequestCtx) error
	Get(key string) interface{}
	GetBytes(key []byte) interface{}
	GetAll() Dict
	Set(key string, value interface{})
	SetBytes(key []byte, value interface{})
	Delete(key string)
	DeleteBytes(key []byte)
	Flush()
	GetSessionID() []byte
	SetSessionID(id []byte)
	Reset()
}

// Provider provider interface
type Provider interface {
	Init(lifeTime int64, cfg ProviderConfig) error
	ReadStore(id []byte) (Storer, error)
	Destroy(id []byte) error
	Regenerate(oldID, newID []byte) (Storer, error)
	Count() int
	NeedGC() bool
	GC()
}

// ProviderConfig provider config interface
type ProviderConfig interface {
	Name() string
}
