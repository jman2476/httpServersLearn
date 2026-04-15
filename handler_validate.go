package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerValidate(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Error decoding chirp", err)
		return
	}

	const maxChirpLength = 140
	currentChirpLength := len(params.Body)
	if currentChirpLength > maxChirpLength {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}

	success := returnVals{Valid: true}
	respondWithJSON(w, 200, success)
	log.Printf("Chirp length good: %d", currentChirpLength)
}
