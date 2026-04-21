package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/jman2476/httpServersLearn/internal/auth"
	"github.com/jman2476/httpServersLearn/internal/database"
)

var (
	errExpiredRefreshToken = errors.New("refresh token: token expired")
	errRevokedRefreshToken = errors.New("refresh token: token revoked")
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	log.Print("POST /api/refresh")

	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "Forbidden", err)
		return
	}

	refreshData, err := cfg.db.GetRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(w, 401, "Forbidden", err)
		return
	}

	if valid, err := cfg.verifyRefreshToken(refreshData); !valid {
		respondWithError(w, 401, "Forbidden", err)
		return
	}

	newToken, err := auth.MakeJWT(refreshData.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, 401, "Forbidden", err)
		return
	}

	respondWithJSON(w, 200, response{Token: newToken})
}

func (cfg *apiConfig) verifyRefreshToken(token database.RefreshToken) (bool, error) {
	if token.RevokedAt.Valid {
		return false, errRevokedRefreshToken
	}

	timeRemaining := time.Until(token.ExpiresAt)

	if timeRemaining <= 0 {
		cfg.db.RevokeToken(context.Background(), token.Token)
		return false, errExpiredRefreshToken
	}

	return true, nil
}
