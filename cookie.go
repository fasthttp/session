package session

import (
	"time"

	"github.com/valyala/fasthttp"
)

// NewCookie return new cookie instance
func NewCookie() *Cookie {
	return new(Cookie)
}

// Get get cookie by name
func (c *Cookie) Get(ctx *fasthttp.RequestCtx, name string) []byte {
	return ctx.Request.Header.Cookie(name)
}

// Set response set cookie
func (c *Cookie) Set(ctx *fasthttp.RequestCtx, name string, value []byte, domain string, expires time.Duration, secure bool) {
	cookie := fasthttp.AcquireCookie()

	cookie.SetKey(name)
	cookie.SetPath("/")
	cookie.SetHTTPOnly(true)
	cookie.SetDomain(domain)
	cookie.SetValueBytes(value)

	if expires >= 0 {
		if expires == 0 {
			cookie.SetExpire(fasthttp.CookieExpireUnlimited)
		} else {
			cookie.SetExpire(time.Now().Add(expires))
		}
	}

	if secure {
		cookie.SetSecure(true)
	}

	ctx.Request.Header.SetCookieBytesKV(cookie.Key(), cookie.Value())
	ctx.Response.Header.SetCookie(cookie)

	fasthttp.ReleaseCookie(cookie)
}

// Delete delete cookie by cookie name
func (c *Cookie) Delete(ctx *fasthttp.RequestCtx, name string) {
	// delete response cookie
	ctx.Response.Header.DelCookie(name)

	// reset response cookie
	cookie := fasthttp.AcquireCookie()

	cookie.SetKey(name)
	cookie.SetValue("")
	cookie.SetPath("/")
	cookie.SetHTTPOnly(true)
	//RFC says 1 second, but let's do it 1 minute to make sure is working...
	exp := time.Now().Add(-time.Duration(1) * time.Minute)
	cookie.SetExpire(exp)
	ctx.Response.Header.SetCookie(cookie)

	// delete request's cookie also
	ctx.Request.Header.DelCookie(name)

	fasthttp.ReleaseCookie(cookie)
}
