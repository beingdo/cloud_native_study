package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type handlerFunction func(http.ResponseWriter, *http.Request)

func unifyHandler(handler handlerFunction) handlerFunction {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				log.Printf("%s %d %s \"%s\"\n", request.RemoteAddr, http.StatusInternalServerError, request.RequestURI, x)
				writer.Write([]byte("server error"))
			}
		}()
		handler(writer, request)
	}
}

func healthzHandler(w http.ResponseWriter, request *http.Request) {
	contentType := request.Header.Get("Content-Type")
	cacheControl := request.Header.Get("Cache-Control")
	connection := request.Header.Get("Connection")
	status := http.StatusOK
	version := os.Getenv("VERSION")
	if version == "" {
		version = "unknown"
	}
	w.Header().Add("Content-Type", contentType)
	w.Header().Add("Cache-Control", cacheControl)
	w.Header().Add("Connection", connection)
	w.Header().Add("Version", version)
	w.WriteHeader(status)
	content := []byte("this is a health check page")
	w.Header().Add("Content-Length", strconv.Itoa(len(content)))
	w.Write(content)
	log.Printf("%s %d %s\n", request.RemoteAddr, status, request.RequestURI)
}
func readinessHandler(w http.ResponseWriter, request *http.Request) {
	contentType := request.Header.Get("Content-Type")
	cacheControl := request.Header.Get("Cache-Control")
	connection := request.Header.Get("Connection")
	status := http.StatusOK
	version := os.Getenv("VERSION")
	if version == "" {
		version = "unknown"
	}
	w.Header().Add("Content-Type", contentType)
	w.Header().Add("Cache-Control", cacheControl)
	w.Header().Add("Connection", connection)
	w.Header().Add("Version", version)
	w.WriteHeader(status)
	content := []byte("this is a health check page")
	w.Header().Add("Content-Length", strconv.Itoa(len(content)))
	w.Write(content)
	log.Printf("%s %d %s\n", request.RemoteAddr, status, request.RequestURI)
}

func panicHandler(w http.ResponseWriter, request *http.Request) {
	contentType := request.Header.Get("Content-Type")
	cacheControl := request.Header.Get("Cache-Control")
	connection := request.Header.Get("Connection")
	status := http.StatusInternalServerError

	w.Header().Add("Content-Type", contentType)
	w.Header().Add("Cache-Control", cacheControl)
	w.Header().Add("Connection", connection)
	w.WriteHeader(status)

	panic("server error")

	log.Printf("%s %d %s\n", request.RemoteAddr, status, request.RequestURI)
}

func defaultHandler(w http.ResponseWriter, request *http.Request) {
	contentType := request.Header.Get("Content-Type")
	cacheControl := request.Header.Get("Cache-Control")
	connection := request.Header.Get("Connection")
	status := http.StatusOK
	version := os.Getenv("VERSION")
	if version == "" {
		version = "unknown"
	}
	w.Header().Add("Content-Type", contentType)
	w.Header().Add("Cache-Control", cacheControl)
	w.Header().Add("Connection", connection)
	w.Header().Add("Version", version)
	w.WriteHeader(status)
	content := []byte("welcome to the default page")
	w.Header().Add("Content-Length", strconv.Itoa(len(content)))
	w.Write(content)
	log.Printf("%s %d %s\n", request.RemoteAddr, status, request.RequestURI)
}

func main() {
	http.HandleFunc("/", unifyHandler(defaultHandler))
	http.HandleFunc("/healthz", unifyHandler(healthzHandler))
	http.HandleFunc("/readiness", unifyHandler(readinessHandler))
	http.HandleFunc("/panic", unifyHandler(panicHandler))
	server := &http.Server{Addr: ":8080"}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		log.Println("server start ...")
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()
	<-signals
	log.Println("SIGTERM received, shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("server shutdown fail:", err)
	}
	log.Println("server shutdown gracefully")

}
