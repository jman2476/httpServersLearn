package main

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"time"

	"github.com/jman2476/httpServersLearn/internal/auth"
	"github.com/jman2476/httpServersLearn/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Error decoding request", err)
		return
	}

	// if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds > 3600 {
	// 	params.ExpiresInSeconds = 3600
	// }

	email, err := mail.ParseAddress(params.Email)
	if err != nil {
		respondWithError(w, 400, "Invalid email address", err)
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(req.Context(), email.Address)
	if err != nil {
		respondWithError(w, 401, "Forbidden", err)
		return
	}

	ok, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if !ok || err != nil {
		respondWithError(w, 401, "Forbidden", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(1)*time.Hour)

	if err != nil {
		respondWithError(w, 500, "Error making token", err)
	}

	refreshArgs := database.CreateRefreshTokenParams{
		Token:  auth.MakeRefreshToken(),
		UserID: user.ID,
	}

	refreshToken, err := cfg.dbQueries.CreateRefreshToken(req.Context(), refreshArgs)
	if err != nil {
		respondWithError(w, 500, "Error making refresh token", err)
	}

	respondWithJSON(w, 200, mapUser(user, token, refreshToken.Token))
}
