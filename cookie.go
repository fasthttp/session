package session

import (
	"time"

	"github.com/valyala/fasthttp"
)

// Cookie cookie struct
type Cookie struct{}

// NewCookie return new cookie instance
func NewCookie() *Cookie {
	return new(Cookie)
}

// Get get cookie by name
func (c *Cookie) Get(ctx *fasthttp.RequestCtx, name string) (value string) {
	cookieByte := ctx.Request.Header.Cookie(name)
	if len(cookieByte) > 0 {
		value = string(cookieByte)
	}
	return
}

// Set response set cookie
func (c *Cookie) Set(ctx *fasthttp.RequestCtx, name string, value string, domain string, expires time.Duration, secure bool) {

	cookie := fasthttp.AcquireCookie()
	defer fasthttp.ReleaseCookie(cookie)

	cookie.SetKey(name)
	cookie.SetPath("/")
	cookie.SetHTTPOnly(true)
	cookie.SetDomain(domain)
	if expires >= 0 {
		if expires == 0 {
			// = 0 unlimited life
			cookie.SetExpire(fasthttp.CookieExpireUnlimited)
		} else {
			// > 0
			cookie.SetExpire(time.Now().Add(expires))
		}
	}
	if ctx.IsTLS() && secure {
		cookie.SetSecure(true)
	}

	cookie.SetValue(value)
	ctx.Response.Header.SetCookie(cookie)
}

// Delete delete cookie by cookie name
func (c *Cookie) Delete(ctx *fasthttp.RequestCtx, name string) {

	// delete response cookie
	ctx.Response.Header.DelCookie(name)

	// reset response cookie
	cookie := fasthttp.AcquireCookie()
	defer fasthttp.ReleaseCookie(cookie)
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
}
