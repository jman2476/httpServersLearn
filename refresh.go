package main

import (
	"net/http"

	"github.com/jman2476/httpServersLearn/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "Forbidden", err)
	}

	refreshData, err := cfg.dbQueries.GetRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(w, 401, "Forbidden", err)
	}
}
