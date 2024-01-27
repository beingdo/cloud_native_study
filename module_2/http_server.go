package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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

func main() {
	http.HandleFunc("/healthz", unifyHandler(healthzHandler))
	http.HandleFunc("/panic", unifyHandler(panicHandler))
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		fmt.Println("server start fail")
	}
}
