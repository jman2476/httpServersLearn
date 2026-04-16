package main

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jman2476/httpServersLearn/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerNewUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
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

	user, err := cfg.dbQueries.CreateUser(req.Context(), email.Address)
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

func mapUser(u database.User) User {
	return User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
	}
}
