package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Forbidden", nil)
		return
	}

	log.Printf("POST /admin/reset")

	cfg.fileserverHits.Store(0)
	err := cfg.dbQueries.ClearUsers(req.Context())
	if err != nil {
		log.Printf("Error clearing users table: %s", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server metrics reset\n"))
	log.Println("Metrics reset")
	log.Println("Users cleared")
}
