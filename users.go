package main

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"

	"github.com/jman2476/httpServersLearn/internal/auth"
	"github.com/jman2476/httpServersLearn/internal/database"
)

func (cfg *apiConfig) handlerNewUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Error decoding email", err)
		return
	}

	email, err := mail.ParseAddress(params.Email)
	if err != nil {
		respondWithError(w, 400, "Invalid email address", err)
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Issue hashing passord", err)
		return
	}

	userArgs := database.CreateUserParams{
		Email:          email.Address,
		HashedPassword: hash,
	}

	user, err := cfg.dbQueries.CreateUser(req.Context(), userArgs)
	if err != nil {
		if strings.Contains(err.Error(), "violates unique constraint \"users_email_key\"") {
			respondWithError(w, 400, "User already exists", err)
			return
		}
		respondWithError(w, 500, "Error creating user", err)
		return
	}

	respondWithJSON(w, 201, mapUser(user))
}
