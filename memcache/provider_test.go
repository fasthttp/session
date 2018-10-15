package memcache

import (
	"testing"

	"github.com/fasthttp/session"
	"github.com/valyala/fasthttp"
)

func getServerSession() *session.Session {
	cfg := session.NewDefaultConfig()
	cfg.SessionIDInHTTPHeader = true // Setted true for simulate the same client in this benchmark
	serverSession := session.New(cfg)
	serverSession.SetProvider(ProviderName, &Config{
		ServerList:   []string{"0.0.0.0"},
		MaxIdleConns: 1,
	})

	return serverSession
}

func Benchmark_Get(b *testing.B) {
	ctx := new(fasthttp.RequestCtx)
	serverSession := getServerSession()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		serverSession.Get(ctx)
	}
}

func Benchmark_Regenerate(b *testing.B) {
	ctx := new(fasthttp.RequestCtx)
	serverSession := getServerSession()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		serverSession.Regenerate(ctx)
	}
}

func Benchmark_Destroy(b *testing.B) {
	ctx := new(fasthttp.RequestCtx)
	serverSession := getServerSession()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		serverSession.Destroy(ctx)
	}
}
