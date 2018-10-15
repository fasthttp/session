package main

import (
	"flag"
	"log"

	"github.com/fasthttp/session"
	"github.com/fasthttp/session/memcache"
	"github.com/fasthttp/session/memory"
	"github.com/fasthttp/session/mysql"
	"github.com/fasthttp/session/postgres"
	"github.com/fasthttp/session/redis"
	"github.com/fasthttp/session/sqlite3"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

const defaultProvider = "memory"

var serverSession = session.New(session.NewDefaultConfig())

func init() {
	providerName := flag.String("provider", defaultProvider, "Name of provider")
	flag.Parse()

	var err error
	switch *providerName {
	case "memory":
		err = serverSession.SetProvider("memory", &memory.Config{})
	case "memcache":
		err = serverSession.SetProvider("memcache", &memcache.Config{
			ServerList: []string{
				"0.0.0.0:11211",
			},
			MaxIdleConns: 8,
			KeyPrefix:    "session",
		})
	case "mysql":
		err = serverSession.SetProvider("mysql", mysql.NewConfigWith("127.0.0.1", 3306, "root", "session", "test", "session"))
	case "postgres":
		err = serverSession.SetProvider("postgres", postgres.NewConfigWith("127.0.0.1", 5432, "root", "session", "test", "session"))
	case "redis":
		err = serverSession.SetProvider("redis", &redis.Config{
			Host:        "127.0.0.1",
			Port:        6379,
			MaxIdle:     8,
			IdleTimeout: 300,
			KeyPrefix:   "session",
		})
	case "sqlite3":
		err = serverSession.SetProvider("sqlite3", sqlite3.NewConfigWith("test.db", "session"))
	default:
		panic("Invalid provider")
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Starting example with provider: " + *providerName)
}

func main() {
	addr := "0.0.0.0:8086"
	router := fasthttprouter.New()
	log.Println("Session example server listen: http://" + addr)

	router.GET("/", indexHandler)
	router.GET("/set", setHandler)
	router.GET("/get", getHandler)
	router.GET("/delete", deleteHandler)
	router.GET("/getAll", getAllHandler)
	router.GET("/flush", flushHandler)
	router.GET("/destroy", destroyHandler)
	router.GET("/sessionid", sessionIdHandler)
	router.GET("/regenerate", regenerateHandler)

	err := fasthttp.ListenAndServe(addr, router.Handler)
	if err != nil {
		log.Println("listen server error :" + err.Error())
	}
}
