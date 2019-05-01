package main

import (
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

// index handler
func indexHandler(ctx *fasthttp.RequestCtx) {
	html := "<h2>Welcome to use session, you should request to the: </h2>"

	html += `> <a href="/">/</a><br>`
	html += `> <a href="/set">set</a><br>`
	html += `> <a href="/get">get</a><br>`
	html += `> <a href="/delete">delete</a><br>`
	html += `> <a href="/getAll">getAll</a><br>`
	html += `> <a href="/flush">flush</a><br>`
	html += `> <a href="/destroy">destroy</a><br>`
	html += `> <a href="/sessionid">sessionid</a><br>`
	html += `> <a href="/regenerate">regenerate</a><br>`

	ctx.SetContentType("text/html;charset=utf-8")
	ctx.SetBodyString(html)
}

// set handler
func setHandler(ctx *fasthttp.RequestCtx) {
	store, err := serverSession.Get(ctx)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	defer serverSession.Save(ctx, store)

	store.Set("foo", "bar")

	ctx.SetBodyString(fmt.Sprintf("Session SET: foo='%s' --> OK", store.Get("foo").(string)))
}

// get handler
func getHandler(ctx *fasthttp.RequestCtx) {
	store, err := serverSession.Get(ctx)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	defer serverSession.Save(ctx, store)

	val := store.Get("foo")
	if val == nil {
		ctx.SetBodyString("Session GET: foo is nil")
		return
	}

	ctx.SetBodyString(fmt.Sprintf("Session GET: foo='%s'", val.(string)))
}

// delete handler
func deleteHandler(ctx *fasthttp.RequestCtx) {
	store, err := serverSession.Get(ctx)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	defer serverSession.Save(ctx, store)

	store.Delete("foo")

	val := store.Get("name")
	if val == nil {
		ctx.SetBodyString("Session DELETE: foo --> OK")
		return
	}
	ctx.SetBodyString("Session DELETE: foo --> ERROR")
}

// get all handler
func getAllHandler(ctx *fasthttp.RequestCtx) {
	store, err := serverSession.Get(ctx)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	defer serverSession.Save(ctx, store)

	store.Set("foo1", "bar1")
	store.Set("foo2", 2)
	store.Set("foo3", "bar3")
	store.Set("foo4", []byte("bar4"))

	data := store.GetAll()

	fmt.Println(data)

	ctx.SetBodyString("Session GetAll: See the OS console!")
}

// flush handle
func flushHandler(ctx *fasthttp.RequestCtx) {
	store, err := serverSession.Get(ctx)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	defer serverSession.Save(ctx, store)

	store.Flush()

	data := store.GetAll()

	fmt.Println(data)

	ctx.SetBodyString("Session FLUSH: See the OS console!")
}

// destroy handle
func destroyHandler(ctx *fasthttp.RequestCtx) {
	err := serverSession.Destroy(ctx)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetBodyString("Session DESTROY --> OK")
}

// get sessionID handle
func sessionIDHandler(ctx *fasthttp.RequestCtx) {
	store, err := serverSession.Get(ctx)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	defer serverSession.Save(ctx, store)

	sessionID := store.GetSessionID()
	ctx.SetBodyString("Session: Current session id: ")
	ctx.Write(sessionID)
}

// regenerate handler
func regenerateHandler(ctx *fasthttp.RequestCtx) {
	store, err := serverSession.Regenerate(ctx)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	defer serverSession.Save(ctx, store)

	sessionID := store.GetSessionID()

	ctx.SetBodyString("Session REGENERATE: New session id: ")
	ctx.Write(sessionID)
}

// get expiration handler
func getExpirationHandler(ctx *fasthttp.RequestCtx) {
	store, err := serverSession.Get(ctx)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	expiration := store.GetExpiration()

	ctx.SetBodyString("Session Expiration: ")
	ctx.WriteString(expiration.String())
}

// set expiration handler
func setExpirationHandler(ctx *fasthttp.RequestCtx) {
	store, err := serverSession.Get(ctx)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	defer serverSession.Save(ctx, store)

	err = store.SetExpiration(30 * time.Second)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetBodyString("Session Expiration set to 30 seconds")
}
