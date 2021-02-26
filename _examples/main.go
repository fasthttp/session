package main

import (
	"flag"
	"log"
	"time"

	"github.com/authelia/session/v2"
	"github.com/authelia/session/v2/providers/memcache"
	"github.com/authelia/session/v2/providers/memory"
	"github.com/authelia/session/v2/providers/mysql"
	"github.com/authelia/session/v2/providers/postgre"
	"github.com/authelia/session/v2/providers/redis"
	"github.com/authelia/session/v2/providers/sqlite3"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

const defaultProvider = "memory"

var serverSession *session.Session

func init() {
	providerName := flag.String("provider", defaultProvider, "Name of provider")
	flag.Parse()

	encoder := session.Base64Encode
	decoder := session.Base64Decode

	var provider session.Provider
	var err error

	switch *providerName {
	case "memory":
		encoder = session.MSGPEncode
		decoder = session.MSGPDecode
		provider, err = memory.New(memory.Config{})
	case "redis":
		encoder = session.MSGPEncode
		decoder = session.MSGPDecode
		provider, err = redis.New(redis.Config{
			KeyPrefix:   "session",
			Addr:        "127.0.0.1:6379",
			PoolSize:    8,
			IdleTimeout: 30 * time.Second,
		})
	case "memcache":
		encoder = session.MSGPEncode
		decoder = session.MSGPDecode
		provider, err = memcache.New(memcache.Config{
			KeyPrefix: "session",
			ServerList: []string{
				"0.0.0.0:11211",
			},
			MaxIdleConns: 8,
		})
	case "mysql":
		cfg := mysql.NewConfigWith("127.0.0.1", 3306, "root", "session", "test", "session")
		provider, err = mysql.New(cfg)
	case "postgre":
		cfg := postgre.NewConfigWith("127.0.0.1", 5432, "postgres", "session", "test", "session")
		provider, err = postgre.New(cfg)
	case "sqlite3":
		cfg := sqlite3.NewConfigWith("test.db", "session")
		provider, err = sqlite3.New(cfg)
	default:
		panic("Invalid provider")
	}

	if err != nil {
		log.Fatal(err)
	}

	cfg := session.NewDefaultConfig()
	cfg.EncodeFunc = encoder
	cfg.DecodeFunc = decoder
	serverSession = session.New(cfg)

	if err = serverSession.SetProvider(provider); err != nil {
		log.Fatal(err)
	}

	log.Print("Starting example with provider: " + *providerName)
}

func main() {
	r := router.New()
	r.GET("/", indexHandler)
	r.GET("/set", setHandler)
	r.GET("/get", getHandler)
	r.GET("/delete", deleteHandler)
	r.GET("/getAll", getAllHandler)
	r.GET("/flush", flushHandler)
	r.GET("/destroy", destroyHandler)
	r.GET("/sessionid", sessionIDHandler)
	r.GET("/regenerate", regenerateHandler)
	r.GET("/setexpiration", setExpirationHandler)
	r.GET("/getexpiration", getExpirationHandler)

	addr := "0.0.0.0:8086"
	log.Println("Session example server listen: http://" + addr)

	err := fasthttp.ListenAndServe(addr, r.Handler)
	if err != nil {
		log.Fatal(err)
	}
}
