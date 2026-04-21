package main

import (
	"log"
	"net/http"

	"github.com/jman2476/httpServersLearn/internal/auth"
)

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, req *http.Request) {
	log.Print("POST /api/revoke")

	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "Forbidden", err)
		return
	}

	err = cfg.dbQueries.RevokeToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(w, 401, "Forbidden", err)
		return
	}

	respondWithJSON(w, 204, struct{}{})
}
