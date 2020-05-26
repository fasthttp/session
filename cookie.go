package session

import (
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

type CookieSameSite int

const (
	// CookieSameSiteDisabled removes the SameSite flag
	CookieSameSiteDisabled CookieSameSite = iota
	// CookieSameSiteDefaultMode sets the SameSite flag
	CookieSameSiteDefaultMode
	// CookieSameSiteLaxMode sets the SameSite flag with the "Lax" parameter
	CookieSameSiteLaxMode
	// CookieSameSiteStrictMode sets the SameSite flag with the "Strict" parameter
	CookieSameSiteStrictMode
	// CookieSameSiteNoneMode sets the SameSite flag with the "None" parameter
	// see https://tools.ietf.org/html/draft-west-cookie-incrementalism-00
	CookieSameSiteNoneMode
)

func newCookie() *cookie {
	return new(cookie)
}

func (c *cookie) get(ctx *fasthttp.RequestCtx, name string) []byte {
	return ctx.Request.Header.Cookie(name)
}

func (c *cookie) set(ctx *fasthttp.RequestCtx, name string, value []byte, domain string, expiration time.Duration, secure bool, sameSite string) {
	cookie := fasthttp.AcquireCookie()

	cookie.SetKey(name)
	cookie.SetPath("/")
	cookie.SetHTTPOnly(true)
	cookie.SetDomain(domain)
	cookie.SetValueBytes(value)

	switch strings.ToLower(sameSite) {
	case "lax":
		cookie.SetSameSite(fasthttp.CookieSameSiteLaxMode)
	case "strict":
		cookie.SetSameSite(fasthttp.CookieSameSiteStrictMode)
	case "none":
		cookie.SetSameSite(fasthttp.CookieSameSiteNoneMode)
	default:
		cookie.SetSameSite(fasthttp.CookieSameSiteDisabled)
	}

	if expiration >= 0 {
		if expiration == 0 {
			cookie.SetExpire(fasthttp.CookieExpireUnlimited)
		} else {
			cookie.SetExpire(time.Now().Add(expiration))
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
