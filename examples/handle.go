package main

import (
	"fmt"

	"github.com/fasthttp/session"
	"github.com/valyala/fasthttp"
)

// request router
func requestRouter(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/":
		indexHandler(ctx)
	case "/set":
		setHandler(ctx)
	case "/get":
		getHandler(ctx)
	case "/delete":
		deleteHandle(ctx)
	case "/getAll":
		getAllHandle(ctx)
	case "/flush":
		flushHandle(ctx)
	case "/destroy":
		destroyHandle(ctx)
	case "/sessionid":
		sessionIDHandle(ctx)
	case "/regenerate":
		regenerateHandle(ctx)
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}

// index handler
func indexHandler(ctx *fasthttp.RequestCtx) {

	html := "<h2>Welcome to use session " + session.Version() + ", you should request to the: </h2>"

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
	// start session
	sessionStore, err := serverSession.Start(ctx)
	if err != nil {
		ctx.SetBodyString(err.Error())
		return
	}
	// must defer sessionStore.save(ctx)
	defer sessionStore.Save(ctx)

	sessionStore.Set("name", "session")

	ctx.SetBodyString(fmt.Sprintf("session setted key name= %s ok", sessionStore.Get("name").(string)))
}

// get handler
func getHandler(ctx *fasthttp.RequestCtx) {
	// start session
	sessionStore, err := serverSession.Start(ctx)
	if err != nil {
		ctx.SetBodyString(err.Error())
		return
	}
	// must defer sessionStore.save(ctx)
	defer sessionStore.Save(ctx)

	s := sessionStore.Get("name")
	if s == nil {
		ctx.SetBodyString("session get name is nil")
		return
	}

	ctx.SetBodyString(fmt.Sprintf("session get name= %s ok", s.(string)))
}

// delete handler
func deleteHandle(ctx *fasthttp.RequestCtx) {
	// start session
	sessionStore, err := serverSession.Start(ctx)
	if err != nil {
		ctx.SetBodyString(err.Error())
		return
	}
	// must defer sessionStore.save(ctx)
	defer sessionStore.Save(ctx)

	sessionStore.Delete("name")

	s := sessionStore.Get("name")
	if s == nil {
		ctx.SetBodyString("session delete key name ok")
		return
	}
	ctx.SetBodyString("session delete key name error")
}

// get all handler
func getAllHandle(ctx *fasthttp.RequestCtx) {
	// start session
	sessionStore, err := serverSession.Start(ctx)
	if err != nil {
		ctx.SetBodyString(err.Error())
		return
	}
	// must defer sessionStore.save(ctx)
	defer sessionStore.Save(ctx)

	sessionStore.Set("foo1", "baa1")
	sessionStore.Set("foo2", "baa2")
	sessionStore.Set("foo3", "baa3")
	sessionStore.Set("foo4", "baa5")

	data := sessionStore.GetAll()

	fmt.Println(data)
	ctx.SetBodyString("session get all data")
}

// flush handle
func flushHandle(ctx *fasthttp.RequestCtx) {
	// start session
	sessionStore, err := serverSession.Start(ctx)
	if err != nil {
		ctx.SetBodyString(err.Error())
		return
	}
	// must defer sessionStore.save(ctx)
	defer sessionStore.Save(ctx)

	sessionStore.Flush()

	ctx.SetBodyString("session flush data")
}

// destroy handle
func destroyHandle(ctx *fasthttp.RequestCtx) {
	// destroy session
	serverSession.Destroy(ctx)

	ctx.SetBodyString("session destroy")
}

// get sessionID handle
func sessionIDHandle(ctx *fasthttp.RequestCtx) {
	// start session
	sessionStore, err := serverSession.Start(ctx)
	if err != nil {
		ctx.SetBodyString(err.Error())
		return
	}
	// must defer sessionStore.save(ctx)
	defer sessionStore.Save(ctx)

	sessionID := sessionStore.GetSessionID()
	ctx.SetBodyString("session sessionID: ")
	ctx.Write(sessionID)
}

// regenerate handler
func regenerateHandle(ctx *fasthttp.RequestCtx) {
	// start session
	sessionStore, err := serverSession.Regenerate(ctx)
	if err != nil {
		ctx.SetBodyString(err.Error())
		return
	}
	// must defer sessionStore.save(ctx)
	defer sessionStore.Save(ctx)

	sessionID := sessionStore.GetSessionID()

	ctx.SetBodyString("session regenerate sessionID: ")
	ctx.Write(sessionID)
}
