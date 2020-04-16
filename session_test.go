package session

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
)

type mockProvider struct {
	errGet        error
	errSave       error
	errDestroy    error
	errRegenerate error
	errGC         error
	countValue    int
	needGCValue   bool
	gcExecuted    bool
}

func (p *mockProvider) Get(store *Store) error {
	return p.errGet
}

func (p *mockProvider) Save(store *Store) error {
	return p.errSave
}

func (p *mockProvider) Destroy(id []byte) error {
	return p.errDestroy
}

func (p *mockProvider) Regenerate(id []byte, newStore *Store) error {
	return p.errRegenerate
}

func (p *mockProvider) Count() int {
	return p.countValue
}

func (p *mockProvider) NeedGC() bool {
	return p.needGCValue
}

func (p *mockProvider) GC() {
	p.gcExecuted = true

	if p.errGC != nil {
		panic(p.errGC)
	}
}

func Test_New(t *testing.T) {
	cfg := Config{
		SessionIDInHTTPHeader: true,
		SessionIDInURLQuery:   true,
	}
	s := New(cfg)

	if s.config.cookieLen != defaultCookieLen {
		t.Errorf("Session.cookieLen == %d, want %d", s.config.cookieLen, defaultCookieLen)
	}

	if s.config.CookieName != defaultSessionKeyName {
		t.Errorf("Session.CookieName == %s, want %s", s.config.CookieName, defaultSessionKeyName)
	}

	if s.config.GCLifetime != defaultGCLifetime {
		t.Errorf("Session.GCLifetime == %d, want %d", s.config.GCLifetime, defaultGCLifetime)
	}

	if s.config.SessionNameInURLQuery != defaultSessionKeyName {
		t.Errorf("Session.SessionNameInURLQuery == %s, want %s", s.config.SessionNameInURLQuery, defaultSessionKeyName)
	}

	if s.config.SessionNameInHTTPHeader != defaultSessionKeyName {
		t.Errorf("Session.SessionNameInHTTPHeader == %s, want %s", s.config.SessionNameInHTTPHeader, defaultSessionKeyName)
	}

	if reflect.ValueOf(s.config.SessionIDGeneratorFunc).Pointer() != reflect.ValueOf(cfg.defaultSessionIDGenerator).Pointer() {
		t.Errorf("Session.SessionIDGeneratorFunc == %p, want %p", s.config.SessionIDGeneratorFunc, cfg.defaultSessionIDGenerator)
	}

	if reflect.ValueOf(s.config.IsSecureFunc).Pointer() != reflect.ValueOf(cfg.defaultIsSecureFunc).Pointer() {
		t.Errorf("Session.IsSecureFunc == %p, want %p", s.config.IsSecureFunc, cfg.defaultIsSecureFunc)
	}

	if s.cookie == nil {
		t.Error("Session.cookie is nil")
	}

	if s.storePool == nil {
		t.Error("Session.storePool is nil")
	}
}

func TestSession_SetProvider(t *testing.T) {
	s := New(Config{
		GCLifetime: 500 * time.Millisecond,
	})
	provider := &mockProvider{needGCValue: true}

	s.SetProvider(provider)
	time.Sleep(s.config.GCLifetime + 100*time.Millisecond)
	s.stopGC()

	if s.provider != provider {
		t.Errorf("Session.SetProvider() provider == %p, want %p", s.provider, provider)
	}

	if !provider.gcExecuted {
		t.Error("GC is not executed")
	}
}

func TestSession_startGC(t *testing.T) {
	s := New(Config{
		GCLifetime: 100 * time.Millisecond,
	})
	provider := &mockProvider{
		needGCValue: true,
		errGC:       errors.New("mock error"),
	}
	s.provider = provider

	defer func() {
		e := recover()
		if e == nil {
			t.Error("Expected panic")
		}
	}()

	s.startGC()

	time.Sleep(s.config.GCLifetime + 100*time.Millisecond)
}

func TestSession_stopGC(t *testing.T) {
	s := New(Config{
		GCLifetime: 100 * time.Millisecond,
	})

	go s.stopGC()

	select {
	case <-s.stopGCChan:
	case <-time.After(200 * time.Millisecond):
		t.Error("Signal for stop GC does not send")
	}
}

func TestSession_setHTTPValues(t *testing.T) {
	ctx := new(fasthttp.RequestCtx)
	s := New(Config{
		SessionIDInHTTPHeader: true,
	})
	id := []byte("sessionID")

	s.setHTTPValues(ctx, id, 100*time.Millisecond)

	if ctx.Response.Header.PeekCookie(s.config.CookieName) == nil {
		t.Error("Session.setHTTPValues() response cookie is not setted")
	}

	if ctx.Request.Header.Cookie(s.config.CookieName) == nil {
		t.Error("Session.setHTTPValues() request cookie is not setted")
	}

	if ctx.Response.Header.Peek(s.config.SessionNameInHTTPHeader) == nil {
		t.Error("Session.setHTTPValues() response header is not setted")
	}

	if ctx.Request.Header.Peek(s.config.SessionNameInHTTPHeader) == nil {
		t.Error("Session.setHTTPValues() request header is not setted")
	}
}

func TestSession_delHTTPValues(t *testing.T) {
	ctx := new(fasthttp.RequestCtx)
	s := New(Config{
		SessionIDInHTTPHeader: true,
	})
	id := []byte("sessionID")

	s.setHTTPValues(ctx, id, 100*time.Millisecond)

	s.delHTTPValues(ctx)

	resultCookie := new(fasthttp.Cookie)
	resultCookie.SetKey(s.config.CookieName)

	if string(resultCookie.Value()) != "" {
		t.Error("Session.setHTTPValues() response cookie is not deleted")
	}

	if ctx.Request.Header.Cookie(s.config.CookieName) != nil {
		t.Error("Session.setHTTPValues() request cookie is not deleted")
	}

	if ctx.Response.Header.Peek(s.config.SessionNameInHTTPHeader) != nil {
		t.Error("Session.setHTTPValues() response header is not deleted")
	}

	if ctx.Request.Header.Peek(s.config.SessionNameInHTTPHeader) != nil {
		t.Error("Session.setHTTPValues() request header is not deleted")
	}
}

func TestSession_getSessionID(t *testing.T) {
	id := "123fvd4r43t4j3tn"

	// From cookie
	s := New(Config{})

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetCookie(s.config.CookieName, id)

	if v := s.getSessionID(ctx); string(v) != id {
		t.Errorf("Session.getSessionID() cookie == %s, want %s", v, id)
	}

	// From header
	s = New(Config{SessionIDInHTTPHeader: true})

	ctx = new(fasthttp.RequestCtx)
	ctx.Request.Header.Set(s.config.SessionNameInHTTPHeader, id)

	if v := s.getSessionID(ctx); string(v) != id {
		t.Errorf("Session.getSessionID() header == %s, want %s", v, id)
	}

	// From url query
	s = New(Config{SessionIDInURLQuery: true})

	ctx = new(fasthttp.RequestCtx)
	ctx.Request.SetRequestURI(fmt.Sprintf("/path?%s=%s", s.config.SessionNameInURLQuery, id))

	if v := s.getSessionID(ctx); string(v) != id {
		t.Errorf("Session.getSessionID() url query == %s, want %s", v, id)
	}
}

func TestSession_GetErrNotProvider(t *testing.T) {
	s := New(Config{})
	ctx := new(fasthttp.RequestCtx)

	store, err := s.Get(ctx)

	if err != errNotSetProvider {
		t.Errorf("Expected error: %v", errNotSetProvider)
	}

	if store != nil {
		t.Error("The store is not nil")
	}
}

func TestSession_GetErrEmptySessionID(t *testing.T) {
	s := New(Config{
		SessionIDGeneratorFunc: func() []byte {
			return []byte("")
		},
	})
	s.SetProvider(new(mockProvider))

	ctx := new(fasthttp.RequestCtx)

	store, err := s.Get(ctx)

	if err != errEmptySessionID {
		t.Errorf("Expected error: %v", errEmptySessionID)
	}

	if store != nil {
		t.Error("The store is not nil")
	}
}

func TestSession_GetProviderError(t *testing.T) {
	s := New(Config{})
	provider := &mockProvider{errGet: errors.New("error from provider")}
	s.SetProvider(provider)

	ctx := new(fasthttp.RequestCtx)

	store, err := s.Get(ctx)

	if err != provider.errGet {
		t.Errorf("Expected error: %v", provider.errGet)
	}

	if store != nil {
		t.Error("The store is not nil")
	}
}

func TestSession_Get(t *testing.T) {
	s := New(Config{})

	provider := new(mockProvider)
	s.SetProvider(provider)

	ctx := new(fasthttp.RequestCtx)

	store, err := s.Get(ctx)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if store == nil {
		t.Error("The store is nil")
	}

	if len(store.sessionID) == 0 {
		t.Error("Store.sessionID is nil")
	}

	if store.defaultExpiration != s.config.Expires {
		t.Errorf("Store.defaultExpiration == %d, want %d", store.defaultExpiration, s.config.Expires)
	}
}

func TestSession_SaveProviderError(t *testing.T) {
	s := New(Config{})
	provider := &mockProvider{errSave: errors.New("error from provider")}
	s.SetProvider(provider)

	ctx := new(fasthttp.RequestCtx)
	store := NewStore()

	err := s.Save(ctx, store)

	if err != provider.errSave {
		t.Errorf("Expected error: %v", provider.errGet)
	}
}

func TestSession_Save(t *testing.T) {
	s := New(Config{})
	provider := new(mockProvider)
	s.SetProvider(provider)

	ctx := new(fasthttp.RequestCtx)

	store, err := s.Get(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if err := s.Save(ctx, store); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if ctx.Response.Header.PeekCookie(s.config.CookieName) == nil {
		t.Error("HTTP values are not setted")
	}

	if len(store.sessionID) > 0 {
		t.Error("store is not reseted")
	}
}

func TestSession_RegenerateErrNotProvider(t *testing.T) {
	s := New(Config{})
	ctx := new(fasthttp.RequestCtx)

	store, err := s.Regenerate(ctx)

	if err != errNotSetProvider {
		t.Errorf("Expected error: %v", errNotSetProvider)
	}

	if store != nil {
		t.Error("The store is not nil")
	}
}

func TestSession_RegenerateErrEmptySessionID(t *testing.T) {
	s := New(Config{
		SessionIDGeneratorFunc: func() []byte {
			return []byte("")
		},
	})
	s.SetProvider(new(mockProvider))

	ctx := new(fasthttp.RequestCtx)

	store, err := s.Regenerate(ctx)

	if err != errEmptySessionID {
		t.Errorf("Expected error: %v", errEmptySessionID)
	}

	if store != nil {
		t.Error("The store is not nil")
	}
}

func TestSession_RegenerateProviderError(t *testing.T) {
	s := New(Config{})
	provider := &mockProvider{errRegenerate: errors.New("error from provider")}
	s.SetProvider(provider)

	ctx := new(fasthttp.RequestCtx)

	store, err := s.Regenerate(ctx)

	if err != provider.errRegenerate {
		t.Errorf("Expected error: %v", provider.errRegenerate)
	}

	if store != nil {
		t.Error("The store is not nil")
	}
}

func TestSession_Regenerate(t *testing.T) {
	s := New(Config{})
	provider := &mockProvider{}
	s.SetProvider(provider)

	ctx := new(fasthttp.RequestCtx)

	store, err := s.Get(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	oldSessionID := string(store.GetSessionID())

	store, err = s.Regenerate(ctx)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if string(store.sessionID) == oldSessionID {
		t.Error("The session id is not chaged")
	}

	if string(ctx.Response.Header.PeekCookie(s.config.CookieName)) == oldSessionID {
		t.Error("HTTP values are not regenerated")
	}
}

func TestSession_DestroyErrNotProvider(t *testing.T) {
	s := New(Config{})
	ctx := new(fasthttp.RequestCtx)

	err := s.Destroy(ctx)

	if err != errNotSetProvider {
		t.Errorf("Expected error: %v", errNotSetProvider)
	}
}

func TestSession_DestroyIDNotExist(t *testing.T) {
	s := New(Config{})
	provider := new(mockProvider)
	s.SetProvider(provider)

	ctx := new(fasthttp.RequestCtx)

	err := s.Destroy(ctx)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestSession_DestroyProviderError(t *testing.T) {
	s := New(Config{})
	provider := &mockProvider{errDestroy: errors.New("error from provider")}
	s.SetProvider(provider)

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetCookie(s.config.CookieName, "asd2324n")

	err := s.Destroy(ctx)

	if err != provider.errDestroy {
		t.Errorf("Expected error: %v", provider.errDestroy)
	}
}

func TestSession_Destroy(t *testing.T) {
	s := New(Config{})
	provider := new(mockProvider)
	s.SetProvider(provider)

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetCookie(s.config.CookieName, "asd2324n")

	err := s.Destroy(ctx)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
