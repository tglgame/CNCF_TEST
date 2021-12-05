package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	logger "github.com/sirupsen/logrus"
)

type configuration struct {
	LogLevel int
	Port     int
	LogPath  string
}

func ReadConfig() *configuration {
	config_file := os.Getenv("CONFIG_FILE_PATH")
	config_file = "conf.json"
	file, _ := os.Open(config_file)
	defer file.Close()

	decorder := json.NewDecoder(file)
	conf := configuration{}
	err := decorder.Decode(&conf)
	if err != nil {
		logger.Error("Error read config file:", err)
	}
	return &conf
}

func Setlog(path string, LogLevel int) {
	writer, _ := rotatelogs.New(
		path+".%Y%m%d%H",
		rotatelogs.WithMaxAge(time.Duration(24*60)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(60)*time.Minute),
	)

	logger.SetOutput(writer)

	if LogLevel == 1 {
		logger.SetLevel(logger.DebugLevel)
	} else if LogLevel == 2 {
		logger.SetLevel(logger.InfoLevel)
	} else if LogLevel == 3 {
		logger.SetLevel(logger.WarnLevel)
	} else if LogLevel == 4 {
		logger.SetLevel(logger.ErrorLevel)
	}
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
	logger.Info("====receive client header====")
	clientIp := GetRemoteIp(req)
	logger.Info("client addr: ", clientIp)

	for k, v := range req.Header {
		logger.Info("receive and set:", k, ":", v)
		w.Header().Set(k, fmt.Sprintf("%s", v))
	}

	var version string
	version = os.Getenv("VERSION")
	w.Header().Set("VERSION", version)
	w.Header().Set("servertag", "servervalue1")

	logger.Info("return status code: ", 200)
	w.WriteHeader(200)

	logger.Info("=====================")
	io.WriteString(w, "hello, world!\n")
}

func Healthz(w http.ResponseWriter, req *http.Request) {
	logger.Info("get health of server")
	w.WriteHeader(200)
}

func main() {
	conf := ReadConfig()
	Setlog(conf.LogPath, conf.LogLevel)
	logger.Info("**********************")
	logger.Info(conf)
	logger.Info(conf.LogPath, conf.LogLevel, conf.Port)
	logger.Info("**********************")

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", HelloServer)
	mux.HandleFunc("/healthz", Healthz)
	srv := http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", conf.Port),
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
	logger.Info("server started!")
	<-done
	logger.Info("server stopped!")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		logger.Info("in cancel context...")
		cancel()
	}()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalln("Server shutdown failed: %+v", err)
	}
	logger.Info("Server gracefully stopped")
}
