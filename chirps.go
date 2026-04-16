package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jman2476/httpServersLearn/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerNewChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Error decoding chirp", err)
		return
	}

	valid, msg := validateChirp(params.Body)
	if !valid {
		respondWithError(w, 400, msg, nil)
		return
	}

	chirpArgs := database.CreateChirpParams{
		Body:   msg,
		UserID: params.UserID,
	}

	chirp, err := cfg.dbQueries.CreateChirp(req.Context(), chirpArgs)
	if err != nil {
		respondWithError(w, 500, "Error creating chirp", err)
		return
	}

	respondWithJSON(w, 201, mapChirp(chirp))
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.dbQueries.GetAllChirps(req.Context())
	if err != nil {
		respondWithError(w, 500, "Error getting chirps", err)
		return
	}

	var allChirps []Chirp
	for _, c := range chirps {
		allChirps = append(allChirps, mapChirp(c))
	}

	respondWithJSON(w, 200, allChirps)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, req *http.Request) {
	chirp_id, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 400, "Invalid UUID value", err)
		return
	}
	chirp, err := cfg.dbQueries.GetChirpByID(req.Context(), chirp_id)
	if err != nil {
		respondWithError(w, 404, "Chirp not found", err)
		return
	}

	respondWithJSON(w, 200, mapChirp(chirp))
}

func mapChirp(c database.Chirp) Chirp {
	return Chirp{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserID:    c.UserID,
	}
}
