package memory

import (
	"testing"

	"github.com/fasthttp/session"
	"github.com/valyala/fasthttp"
)

func getServerSession() *session.Session {
	cfg := session.NewDefaultConfig()
	cfg.SessionIDInHTTPHeader = true // Setted true for simulate the same client in this benchmark
	serverSession := session.New(cfg)
	serverSession.SetProvider(ProviderName, &Config{})

	return serverSession
}

func Benchmark_Get(b *testing.B) {
	testCtx := new(fasthttp.RequestCtx)
	serverSession := getServerSession()

	handler := func(ctx *fasthttp.RequestCtx) {
		store, _ := serverSession.Get(ctx)
		store.Set("k1", 1)
		serverSession.Save(store)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler(testCtx)
	}
}

func Benchmark_Regenerate(b *testing.B) {
	testCtx := new(fasthttp.RequestCtx)
	serverSession := getServerSession()

	handler := func(ctx *fasthttp.RequestCtx) {
		store, _ := serverSession.Regenerate(ctx)
		store.Set("k1", 1)
		serverSession.Save(store)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler(testCtx)
	}
}

func Benchmark_Destroy(b *testing.B) {
	testCtx := new(fasthttp.RequestCtx)
	serverSession := getServerSession()

	handler := func(ctx *fasthttp.RequestCtx) {
		serverSession.Destroy(ctx)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler(testCtx)
	}
}
