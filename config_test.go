package session

import (
	"reflect"
	"testing"

	"github.com/valyala/fasthttp"
)

func Test_NewDefaultConfig(t *testing.T) {
	cfg := NewDefaultConfig()

	if cfg.CookieName != defaultSessionKeyName {
		t.Errorf("NewDefaultConfig().CookieName == %s, want %s", cfg.CookieName, defaultSessionKeyName)
	}

	if cfg.Domain != defaultDomain {
		t.Errorf("NewDefaultConfig().Domain == %s, want %s", cfg.Domain, defaultDomain)
	}

	if cfg.Expiration != defaultExpiration {
		t.Errorf("NewDefaultConfig().Expiration == %d, want %d", cfg.Expiration, defaultExpiration)
	}

	if cfg.GCLifetime != defaultGCLifetime {
		t.Errorf("NewDefaultConfig().GCLifetime == %d, want %d", cfg.GCLifetime, defaultGCLifetime)
	}

	if cfg.Secure != defaultSecure {
		t.Errorf("NewDefaultConfig().Secure == %v, want %v", cfg.Secure, defaultSecure)
	}

	if cfg.SessionIDInURLQuery != defaultSessionIDInURLQuery {
		t.Errorf("NewDefaultConfig().SessionIDInURLQuery == %v, want %v", cfg.SessionIDInURLQuery, defaultSessionIDInURLQuery)
	}

	if cfg.SessionNameInURLQuery != defaultSessionKeyName {
		t.Errorf("NewDefaultConfig().SessionNameInURLQuery == %s, want %s", cfg.SessionNameInURLQuery, defaultSessionKeyName)
	}

	if cfg.SessionIDInHTTPHeader != defaultSessionIDInHTTPHeader {
		t.Errorf("NewDefaultConfig().SessionIDInHTTPHeader == %v, want %v", cfg.SessionIDInHTTPHeader, defaultSessionIDInHTTPHeader)
	}

	if cfg.SessionNameInHTTPHeader != defaultSessionKeyName {
		t.Errorf("NewDefaultConfig().SessionNameInHTTPHeader == %s, want %s", cfg.SessionNameInHTTPHeader, defaultSessionKeyName)
	}

	if cfg.cookieLen != defaultCookieLen {
		t.Errorf("NewDefaultConfig().cookieLen == %d, want %d", cfg.cookieLen, defaultCookieLen)
	}

	if reflect.ValueOf(cfg.SessionIDGeneratorFunc).Pointer() != reflect.ValueOf(cfg.defaultSessionIDGenerator).Pointer() {
		t.Errorf("NewDefaultConfig().SessionIDGeneratorFunc == %p, want %p", cfg.SessionIDGeneratorFunc, cfg.defaultSessionIDGenerator)
	}

	if reflect.ValueOf(cfg.IsSecureFunc).Pointer() != reflect.ValueOf(cfg.defaultIsSecureFunc).Pointer() {
		t.Errorf("NewDefaultConfig().IsSecureFunc == %p, want %p", cfg.IsSecureFunc, cfg.defaultIsSecureFunc)
	}
}

func TestConfig_defaultSessionIDGenerator(t *testing.T) {
	cfg := Config{cookieLen: 10}

	id := cfg.defaultSessionIDGenerator()

	if len(id) != int(cfg.cookieLen) {
		t.Errorf("Config.defaultSessionIDGenerator() len == %d, want %d", len(id), cfg.cookieLen)
	}
}

func TestConfig_defaultIsSecureFunc(t *testing.T) {
	cfg := Config{}
	ctx := new(fasthttp.RequestCtx)

	secure := cfg.defaultIsSecureFunc(ctx)

	if secure != ctx.IsTLS() {
		t.Errorf("Config.defaultIsSecureFunc() == %v, want %v", secure, ctx.IsTLS())
	}
}
