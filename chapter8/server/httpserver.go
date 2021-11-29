package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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
		logger.Println("receive and set:", k, ":", v)
		w.Header().Set(k, fmt.Sprintf("%s", v))
	}

	var version string
	version = os.Getenv("VERSION")
	w.Header().Set("VERSION", version)
	w.Header().Set("servertag", "servervalue1")

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

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", HelloServer)
	mux.HandleFunc("/healthz", Healthz)
	srv := http.Server{
		Addr:    "0.0.0.0:80",
		Handler: mux,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// server start
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen and serve: ", err)
		}
	}()
	logger.Println("server started!")
	<-done
	logger.Println("server stopped!")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		logger.Println("in cancel context...")
		cancel()
	}()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalln("Server shutdown failed: %+v", err)
	}
	logger.Println("Server gracefully stopped")
}
