package main

import (
	"log"
	"net/http"
)

func (c *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	c.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server metrics reset\n"))
	log.Println("Metrics reset")
}
