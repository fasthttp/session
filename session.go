package session

import (
	"errors"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

var providers Dict

// Register register session provider
func Register(providerName string, provider Provider) {
	if provider == nil {
		panic("session register error, provider " + providerName + " is nil!")
	} else if providers.Has(providerName) {
		panic("session register error, provider " + providerName + " already registered!")
	}

	providers.Set(providerName, provider)
}

// New return new Session
func New(cfg *Config) *Session {
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
	if cfg.SessionLifetime == 0 {
		cfg.SessionLifetime = defaultSessionLifetime
	}

	if cfg.SessionIDGeneratorFunc == nil {
		cfg.SessionIDGeneratorFunc = cfg.defaultSessionIDGenerator
	}

	session := &Session{
		config: cfg,
		cookie: NewCookie(),
	}

	return session
}

// SetProvider set session provider and provider config
func (s *Session) SetProvider(name string, cfg ProviderConfig) error {
	if !providers.Has(name) {
		return errors.New("session set provider error, " + name + " not registered!")
	}
	s.provider = providers.Get(name).(Provider)

	err := s.provider.Init(s.config.SessionLifetime, cfg)
	if err != nil {
		return err
	}

	if s.provider.NeedGC() {
		s.startGC()
	}

	return nil
}

// startGC start session gc process.
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

func (s *Session) setHTTPValues(ctx *fasthttp.RequestCtx, sessionID []byte) {
	s.cookie.Set(ctx, s.config.CookieName, sessionID, s.config.Domain, s.config.Expires, s.config.Secure)

	if s.config.SessionIDInHTTPHeader {
		ctx.Request.Header.SetBytesV(s.config.SessionNameInHTTPHeader, sessionID)
		ctx.Response.Header.SetBytesV(s.config.SessionNameInHTTPHeader, sessionID)
	}
}

func (s *Session) delHTTPValues(ctx *fasthttp.RequestCtx) {
	s.cookie.Delete(ctx, s.config.CookieName)

	if s.config.SessionIDInHTTPHeader {
		ctx.Request.Header.Del(s.config.SessionNameInHTTPHeader)
		ctx.Response.Header.Del(s.config.SessionNameInHTTPHeader)
	}
}

// get session id
// 1. get session id from cookie
// 2. get session id from http headers
// 3. get session id from query string
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

// Get get user session from provider
// 1. get sessionID from fasthttp ctx
// 2. if sessionID is empty, generator sessionID and set response Set-Cookie
// 3. return session provider store
func (s *Session) Get(ctx *fasthttp.RequestCtx) (Storer, error) {
	if s.provider == nil {
		return nil, errNotSetProvider
	}

	sessionID := s.getSessionID(ctx)
	if len(sessionID) == 0 {
		sessionID = s.config.SessionIDGeneratorFunc()
		if len(sessionID) == 0 {
			return nil, errEmptySessionID
		}

		s.setHTTPValues(ctx, sessionID)
	}

	store, err := s.provider.Get(sessionID)
	if err != nil {
		return nil, err
	}

	return store, nil
}

// Save save user session with current store
func (s *Session) Save(store Storer) error {
	err := store.Save()
	if err != nil {
		return err
	}

	fmt.Printf("SAVE --- %p\n", store)

	s.provider.Put(store)

	return nil
}

// Regenerate regenerate a session id for this Storer
func (s *Session) Regenerate(ctx *fasthttp.RequestCtx) (Storer, error) {
	if s.provider == nil {
		return nil, errNotSetProvider
	}

	newID := s.config.SessionIDGeneratorFunc()
	if len(newID) == 0 {
		return nil, errEmptySessionID
	}
	oldID := s.getSessionID(ctx)

	store, err := s.provider.Regenerate(oldID, newID)
	if err != nil {
		return nil, err
	}

	s.setHTTPValues(ctx, newID)

	return store, nil
}

// Destroy destroy session in fasthttp ctx
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

// Version return session version
func Version() string {
	return version
}
