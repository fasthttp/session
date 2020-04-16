package session

import (
	"time"

	"github.com/valyala/fasthttp"
)

func newCookie() *cookie {
	return new(cookie)
}

func (c *cookie) get(ctx *fasthttp.RequestCtx, name string) []byte {
	return ctx.Request.Header.Cookie(name)
}

func (c *cookie) set(ctx *fasthttp.RequestCtx, name string, value []byte, domain string, expires time.Duration, secure bool) {
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

func (c *cookie) delete(ctx *fasthttp.RequestCtx, name string) {
	ctx.Request.Header.DelCookie(name)
	ctx.Response.Header.DelCookie(name)

	cookie := fasthttp.AcquireCookie()
	cookie.SetKey(name)
	cookie.SetValue("")
	cookie.SetPath("/")
	cookie.SetHTTPOnly(true)
	//RFC says 1 second, but let's do it 1 minute to make sure is working...
	exp := time.Now().Add(-1 * time.Minute)
	cookie.SetExpire(exp)
	ctx.Response.Header.SetCookie(cookie)

	fasthttp.ReleaseCookie(cookie)
}
