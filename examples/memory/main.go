package main

// session memory provider example

import (
	"log"
	"os"

	"github.com/fasthttp/session"
	"github.com/fasthttp/session/memory"
	"github.com/valyala/fasthttp"
)

// default config
var serverSession = session.NewSession(session.NewDefaultConfig())

// custom config
//var serverSession = session.NewSession(&session.Config{
//	CookieName: "ssid",
//	Domain: "",
//	Expires: time.Hour * 2,
//	GCLifetime: 3,
//	SessionLifetime: 60,
//	Secure: true,
//	SessionIDInURLQuery: false,
//	SessionNameInURLQuery: "",
//	SessionIDInHTTPHeader: false,
//	SessionNameInHTTPHeader: "",
//	SessionIDGeneratorFunc: func() string {return ""},
//	EncodeFunc: func(cookieValue string) (string, error) {return "", nil},
//	DecodeFunc: func(cookieValue string) (string, error) {return "", nil},
//})

func main() {

	// You must set up provider before use
	err := serverSession.SetProvider("memory", &memory.Config{})
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	addr := ":8086"
	log.Println("session memory example server listen: " + addr)
	// Fasthttp start listen serve
	err = fasthttp.ListenAndServe(addr, requestRouter)
	if err != nil {
		log.Println("listen server error :" + err.Error())
	}
}
