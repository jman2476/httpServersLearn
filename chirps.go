package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/jman2476/httpServersLearn/internal/auth"
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
		Body string `json:"body"`
	}

	log.Printf("POST /api/chirps")

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "Forbidden", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "Forbidden", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
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
		UserID: userID,
	}

	chirp, err := cfg.db.CreateChirp(req.Context(), chirpArgs)
	if err != nil {
		respondWithError(w, 500, "Error creating chirp", err)
		return
	}

	respondWithJSON(w, 201, mapChirp(chirp))
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, req *http.Request) {
	log.Printf("GET /api/chirps")

	authorID := req.URL.Query().Get("author_id")

	if authorID != "" {
		cfg.handlerGetChirpsByAuthor(w, req, authorID)
		return
	}

	chirps, err := cfg.db.GetAllChirps(req.Context())
	if err != nil {
		respondWithError(w, 500, "Error getting chirps", err)
		return
	}

	var allChirps []Chirp
	for _, c := range chirps {
		allChirps = append(allChirps, mapChirp(c))
	}

	sortOrder := req.URL.Query().Get("sort")
	if sortOrder == "desc" {
		slices.SortFunc(allChirps, func(a, b Chirp) int {
			return b.CreatedAt.Compare(a.CreatedAt)
		})
	}

	respondWithJSON(w, 200, allChirps)
}

func (cfg *apiConfig) handlerGetChirpsByAuthor(w http.ResponseWriter, req *http.Request, authorID string) {
	id, err := uuid.Parse(authorID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Malformed author id", err)
		return
	}

	chirps, err := cfg.db.GetChirpsByAuthor(req.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No entries found", err)
		return
	}

	var authorChirps []Chirp
	for _, c := range chirps {
		authorChirps = append(authorChirps, mapChirp(c))
	}

	sortOrder := req.URL.Query().Get("sort")
	if sortOrder == "desc" {
		slices.SortFunc(authorChirps, func(a, b Chirp) int {
			return b.CreatedAt.Compare(a.CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, authorChirps)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, req *http.Request) {
	log.Printf("GET /api/chirps/{chirpID}")

	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 400, "Invalid UUID value", err)
		return
	}
	chirp, err := cfg.db.GetChirpByID(req.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "Chirp not found", err)
		return
	}

	respondWithJSON(w, 200, mapChirp(chirp))
}

func (cfg *apiConfig) handlerDeleteChirpByID(w http.ResponseWriter, req *http.Request) {
	log.Print("DELETE /api/chirps/{chirpID}")

	authToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Bad token", err)
		return
	}

	userID, err := auth.ValidateJWT(authToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Bad token", err)
		return
	}

	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	chirp, err := cfg.db.GetChirpByID(req.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Not authorized", err)
		return
	}

	err = cfg.db.DeleteChirpByID(req.Context(), chirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})

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
