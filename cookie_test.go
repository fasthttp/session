package session

import (
	"testing"
	"time"

	"github.com/valyala/fasthttp"
)

func Test_newCookie(t *testing.T) {
	cookie := newCookie()

	if cookie == nil {
		t.Error("newCookie() return nil")
	}
}

func TestCookie_get(t *testing.T) {
	ctx := new(fasthttp.RequestCtx)
	cookie := newCookie()

	key := "key"
	value := "value"

	ctx.Request.Header.SetCookie(key, value)

	if v := cookie.get(ctx, key); string(v) != value {
		t.Errorf("cookie.get() == %s, want %s", v, value)
	}
}

func TestCookie_set(t *testing.T) {
	ctx := new(fasthttp.RequestCtx)
	cookie := newCookie()

	key := "key"
	value := []byte("value")
	path := "/"
	domain := "domain"
	expiration := 10 * time.Second
	secure := true
	samesite := "Lax"

	now := time.Now()
	cookie.set(ctx, key, value, domain, expiration, secure, samesite)

	resultCookie := new(fasthttp.Cookie)
	resultCookie.SetKey(key)
	ctx.Response.Header.Cookie(resultCookie)

	if string(resultCookie.Path()) != path {
		t.Errorf("cookie.set() Path == %s, want %s", resultCookie.Path(), path)
	}

	if !resultCookie.HTTPOnly() {
		t.Errorf("cookie.set() HTTPOnly == %v, want %v", false, true)
	}

	if string(resultCookie.Domain()) != domain {
		t.Errorf("cookie.set() Domain == %s, want %s", resultCookie.Domain(), domain)
	}

	if string(resultCookie.Value()) != string(value) {
		t.Errorf("cookie.set() Value == %s, want %s", resultCookie.Value(), value)
	}

	if resultCookie.Expire().Unix() != now.Add(expiration).Unix() {
		t.Errorf("cookie.set() Expire == %v, want %v", resultCookie.Expire(), expiration)
	}

	if resultCookie.Secure() != secure {
		t.Errorf("cookie.set() Secure == %v, want %v", resultCookie.Secure(), secure)
	}

	if v := ctx.Request.Header.Cookie(key); string(v) != string(value) {
		t.Errorf("cookie.set() request value == %s, want %s", v, value)
	}
}

func TestCookie_delete(t *testing.T) {
	ctx := new(fasthttp.RequestCtx)
	cookie := newCookie()

	key := "key"
	value := []byte("")
	path := "/"
	expiration := -1 * time.Minute

	now := time.Now()
	cookie.delete(ctx, key)

	resultCookie := new(fasthttp.Cookie)
	resultCookie.SetKey(key)
	ctx.Response.Header.Cookie(resultCookie)

	if string(resultCookie.Path()) != path {
		t.Errorf("cookie.set() Path == %s, want %s", resultCookie.Path(), path)
	}

	if !resultCookie.HTTPOnly() {
		t.Errorf("cookie.set() HTTPOnly == %v, want %v", false, true)
	}

	if string(resultCookie.Value()) != string(value) {
		t.Errorf("cookie.set() Value == %s, want %s", resultCookie.Value(), value)
	}

	if resultCookie.Expire().Unix() != now.Add(expiration).Unix() {
		t.Errorf("cookie.set() Expire == %v, want %v", resultCookie.Expire(), expiration)
	}

	if v := ctx.Request.Header.Cookie(key); string(v) != string(value) {
		t.Errorf("cookie.set() request value == %s, want %s", v, value)
	}
}
