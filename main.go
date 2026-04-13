package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filepathRoot = "."

	serveMux := http.NewServeMux() //server multiplexer

	serveMux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	log.Printf("Start serving on port %s", port)
	log.Fatal(server.ListenAndServe())
}
