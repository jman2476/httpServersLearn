package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"strings"

	"github.com/jman2476/httpServersLearn/internal/auth"
	"github.com/jman2476/httpServersLearn/internal/database"
)

func (cfg *apiConfig) handlerNewUser(w http.ResponseWriter, req *http.Request) {
	log.Printf("POST /api/users")

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

	user, err := cfg.db.CreateUser(req.Context(), userArgs)
	if err != nil {
		if strings.Contains(err.Error(), "violates unique constraint \"users_email_key\"") {
			respondWithError(w, 400, "User already exists", err)
			return
		}
		respondWithError(w, 500, "Error creating user", err)
		return
	}

	respondWithJSON(w, 201, mapUser(user, "", ""))
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, req *http.Request) {
	log.Print("PUT /api/users")

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	authToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Bad token", err)
		return
	}

	userID, err := auth.ValidateJWT(authToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Bad token", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Need email and password", err)
		return
	}

	hashPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal error", err)
		return
	}

	updateArgs := database.UpdateUserByIDParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashPassword,
	}

	user, err := cfg.db.UpdateUserByID(req.Context(), updateArgs)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error", err)
		return
	}

	respondWithJSON(w, http.StatusOK, mapUser(user, "", ""))
}
