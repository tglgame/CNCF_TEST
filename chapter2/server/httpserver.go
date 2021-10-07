package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
)

var logger *log.Logger
var once sync.Once

func GetLogInstance() *log.Logger {
	once.Do(func() {
		logger = log.New(os.Stdout, "http_server:", log.Ldate|log.Ltime|log.Lshortfile)
	})
	return logger
}

func GetRemoteIp(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get(`X-Real-Ip`); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get(`X-Forwarded-For`); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}
	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	logger := GetLogInstance()

	logger.Println("====receive client header====")

	clientIp := GetRemoteIp(req)
	logger.Println("client addr: ", clientIp)

	for k, v := range req.Header {
		logger.Println("receive and set:", k, ":", v[0])
		w.Header().Set(k, v[0])
	}

	var version string
	version = os.Getenv("VERSION")
	w.Header().Set("VERSION", version)

	w.Header().Set("server1", "sv1")
	w.Header().Set("server2", "sv2")

	logger.Println("return status code: ", 200)
	w.WriteHeader(200)

	logger.Println("=====================")
	io.WriteString(w, "hello, world!\n")
}

func Healthz(w http.ResponseWriter, req *http.Request) {
	logger := GetLogInstance()
	logger.Println("get health of server")
	w.WriteHeader(200)
}

func main() {
	logger := GetLogInstance()

	http.HandleFunc("/hello", HelloServer)
	http.HandleFunc("/healthz", Healthz)
	err := http.ListenAndServe("0.0.0.0:80", nil)
	if err != nil {
		logger.Fatal("listen and serve: ", err)
	}
}
