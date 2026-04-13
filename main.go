package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filepathRoot = "."

	serveMux := http.NewServeMux() //server multiplexer

	rmPrefixHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	serveMux.Handle("/app/", rmPrefixHandler)
	serveMux.HandleFunc("/healthz", handlerHealthz)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	log.Printf("Start serving from %s on port %s", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func handlerHealthz(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	// w.Header().Set("charset", "utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
