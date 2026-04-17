package main

import (
	"encoding/json"
	"net/http"
	"net/mail"

	"github.com/jman2476/httpServersLearn/internal/auth"
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

	respondWithJSON(w, 200, mapUser(user))
}
