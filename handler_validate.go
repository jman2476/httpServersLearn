package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func handlerValidate(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Body string `json:"cleaned_body"`
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

	respondWithJSON(w, 200, returnVals{
		Body: cleanChirp(params.Body),
	})
	log.Printf("Chirp length good: %d", currentChirpLength)
}

func cleanChirp(chirp string) string {
	profanity := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	var cleaned strings.Builder
	words := strings.Split(chirp, " ")
	var cleanWords []string

	for _, word := range words {
		if slices.Contains(profanity, strings.ToLower(word)) {
			log.Println("Profanity found")
			word = "****"
		}
		cleanWords = append(cleanWords, word)
	}
	cleaned.WriteString(strings.Join(cleanWords, " "))

	return cleaned.String()
}
