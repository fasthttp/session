package session

import (
	"fmt"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

// New returns a configured manager
func New(cfg Config) *Session {
	cfg.cookieLen = defaultCookieLen

	if cfg.CookieName == "" {
		cfg.CookieName = defaultSessionKeyName
	}
	if cfg.SessionIDInHTTPHeader && cfg.SessionNameInHTTPHeader == "" {
		cfg.SessionNameInHTTPHeader = defaultSessionKeyName
	}
	if cfg.SessionIDInURLQuery && cfg.SessionNameInURLQuery == "" {
		cfg.SessionNameInURLQuery = defaultSessionKeyName
	}

	if cfg.GCLifetime == 0 {
		cfg.GCLifetime = defaultGCLifetime
	}

	if cfg.SessionIDGeneratorFunc == nil {
		cfg.SessionIDGeneratorFunc = cfg.defaultSessionIDGenerator
	}

	if cfg.IsSecureFunc == nil {
		cfg.IsSecureFunc = cfg.defaultIsSecureFunc
	}

	session := &Session{
		config: cfg,
		cookie: newCookie(),
		storePool: &sync.Pool{
			New: func() interface{} {
				return NewStore()
			},
		},
	}

	return session
}

// SetProvider sets the session provider used by the sessions manager
func (s *Session) SetProvider(provider Provider) error {
	s.provider = provider

	if s.provider.NeedGC() {
		s.startGC()
	}

	return nil
}

func (s *Session) startGC() {
	go func() {
		defer func() {
			e := recover()
			if e != nil {
				panic(fmt.Errorf("session gc crash, %v", e))
			}
		}()

		for {
			select {
			case <-time.After(s.config.GCLifetime):
				s.provider.GC()
			}
		}
	}()
}

func (s *Session) setHTTPValues(ctx *fasthttp.RequestCtx, sessionID []byte, expires time.Duration) {
	secure := s.config.Secure && s.config.IsSecureFunc(ctx)
	s.cookie.set(ctx, s.config.CookieName, sessionID, s.config.Domain, expires, secure)

	if s.config.SessionIDInHTTPHeader {
		ctx.Request.Header.SetBytesV(s.config.SessionNameInHTTPHeader, sessionID)
		ctx.Response.Header.SetBytesV(s.config.SessionNameInHTTPHeader, sessionID)
	}
}

func (s *Session) delHTTPValues(ctx *fasthttp.RequestCtx) {
	s.cookie.delete(ctx, s.config.CookieName)

	if s.config.SessionIDInHTTPHeader {
		ctx.Request.Header.Del(s.config.SessionNameInHTTPHeader)
		ctx.Response.Header.Del(s.config.SessionNameInHTTPHeader)
	}
}

// Returns the session id from:
// 1. cookie
// 2. http headers
// 3. query string
func (s *Session) getSessionID(ctx *fasthttp.RequestCtx) []byte {
	val := ctx.Request.Header.Cookie(s.config.CookieName)
	if len(val) > 0 {
		return val
	}

	if s.config.SessionIDInHTTPHeader {
		val = ctx.Request.Header.Peek(s.config.SessionNameInHTTPHeader)
		if len(val) > 0 {
			return val
		}
	}

	if s.config.SessionIDInURLQuery {
		val = ctx.FormValue(s.config.SessionNameInURLQuery)
		if len(val) > 0 {
			return val
		}
	}

	return nil
}

// Get returns the user session
// if it does not exist, it will be generated
func (s *Session) Get(ctx *fasthttp.RequestCtx) (*Store, error) {
	if s.provider == nil {
		return nil, errNotSetProvider
	}

	id := s.getSessionID(ctx)
	if len(id) == 0 {
		id = s.config.SessionIDGeneratorFunc()
		if len(id) == 0 {
			return nil, errEmptySessionID
		}
	}

	store := s.storePool.Get().(*Store)
	store.sessionID = id
	store.defaultExpiration = s.config.Expires

	if err := s.provider.Get(store); err != nil {
		return nil, err
	}

	return store, nil
}

// Save saves the user session
//
// Warning: Don't use the store after exec this function, because, you will lose the after data
// For avoid it, defer this function in your request handler
func (s *Session) Save(ctx *fasthttp.RequestCtx, store *Store) error {
	if err := s.provider.Save(store); err != nil {
		return err
	}

	s.setHTTPValues(ctx, store.GetSessionID(), store.GetExpiration())

	store.Reset()
	s.storePool.Put(store)

	return nil
}

// Regenerate generates a new session id to the current user
func (s *Session) Regenerate(ctx *fasthttp.RequestCtx) (*Store, error) {
	if s.provider == nil {
		return nil, errNotSetProvider
	}

	newID := s.config.SessionIDGeneratorFunc()
	if len(newID) == 0 {
		return nil, errEmptySessionID
	}
	id := s.getSessionID(ctx)

	store := s.storePool.Get().(*Store)
	store.sessionID = newID
	store.defaultExpiration = s.config.Expires

	if err := s.provider.Regenerate(id, store); err != nil {
		return nil, err
	}

	s.setHTTPValues(ctx, newID, store.GetExpiration())

	return store, nil
}

// Destroy destroys the session of the current user
func (s *Session) Destroy(ctx *fasthttp.RequestCtx) error {
	sessionID := s.getSessionID(ctx)
	if len(sessionID) == 0 {
		return nil
	}

	err := s.provider.Destroy(sessionID)
	if err != nil {
		return err
	}

	s.delHTTPValues(ctx)

	return nil
}
