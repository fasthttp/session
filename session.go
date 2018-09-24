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
		cfg.SessionLifetime = cfg.GCLifetime
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
	provider := providers.Get(name).(Provider)
	if provider == nil {
		return errors.New("session set provider error, " + name + " not registered!")
	}

	err := provider.Init(s.config.SessionLifetime, cfg)
	if err != nil {
		return err
	}

	s.provider = provider

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

		timeAfter := time.Duration(s.config.GCLifetime) * time.Second
		for {
			select {
			case <-time.After(timeAfter):
				s.provider.GC()
			}
		}
	}()
}

func (s *Session) setHTTPValues(ctx *fasthttp.RequestCtx, encodeCookieValue, sessionID []byte) {
	s.cookie.Set(ctx, s.config.CookieName, encodeCookieValue, s.config.Domain, s.config.Expires, s.config.Secure)

	if s.config.SessionIDInHTTPHeader {
		ctx.Request.Header.SetBytesV(s.config.SessionNameInHTTPHeader, sessionID)
		ctx.Response.Header.SetBytesV(s.config.SessionNameInHTTPHeader, sessionID)
	}
}

// GetSessionID get session id
// 1. get session id by reading from cookie
// 2. get session id from query
// 3. get session id from http headers
func (s *Session) GetSessionID(ctx *fasthttp.RequestCtx) []byte {
	val := ctx.Request.Header.Cookie(s.config.CookieName)
	if len(val) > 0 {
		return s.config.Decode(val)
	}

	if s.config.SessionIDInURLQuery {
		val = ctx.FormValue(s.config.SessionNameInURLQuery)
		if len(val) > 0 {
			return s.config.Decode(val)
		}

	}
	if s.config.SessionIDInHTTPHeader {
		val = ctx.Request.Header.Peek(s.config.SessionNameInHTTPHeader)
		if len(val) > 0 {
			return s.config.Decode(val)
		}
	}

	return nil
}

// Start session start
// 1. get sessionID from fasthttp ctx
// 2. if sessionID is empty, generator sessionID and set response Set-Cookie
// 3. return session provider store
func (s *Session) Start(ctx *fasthttp.RequestCtx) (Storer, error) {
	if s.provider == nil {
		return nil, errNotSetProvider
	}

	sessionID := s.GetSessionID(ctx)
	if len(sessionID) == 0 {
		sessionID = s.config.SessionIDGenerator()
		if len(sessionID) == 0 {
			return nil, errEmptySessionID
		}

		s.setHTTPValues(ctx, s.config.Encode(sessionID), sessionID)
	}

	store, err := s.provider.ReadStore(sessionID)
	if err != nil {
		return nil, err
	}

	return store, nil
}

// Regenerate regenerate a session id for this Storer
func (s *Session) Regenerate(ctx *fasthttp.RequestCtx) (Storer, error) {
	if s.provider == nil {
		return nil, errNotSetProvider
	}

	newSessionID := s.config.SessionIDGenerator()
	if len(newSessionID) == 0 {
		return nil, errEmptySessionID
	}
	oldSessionID := s.GetSessionID(ctx)

	store, err := s.provider.Regenerate(oldSessionID, newSessionID)
	if err != nil {
		return nil, err
	}

	s.setHTTPValues(ctx, s.config.Encode(newSessionID), newSessionID)

	return store, nil
}

// Destroy destroy session in fasthttp ctx
func (s *Session) Destroy(ctx *fasthttp.RequestCtx) {
	sessionID := s.GetSessionID(ctx)
	if len(sessionID) == 0 {
		return
	}

	s.provider.Destroy(sessionID)
	s.cookie.Delete(ctx, s.config.CookieName)

	if s.config.SessionIDInHTTPHeader {
		ctx.Request.Header.Del(s.config.SessionNameInHTTPHeader)
		ctx.Response.Header.Del(s.config.SessionNameInHTTPHeader)
	}
}

// Version return session version
func Version() string {
	return version
}
