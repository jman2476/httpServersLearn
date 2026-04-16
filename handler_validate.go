package main

import (
	"log"
	"slices"
	"strings"
)

func validateChirp(chirp string) (bool, string) {
	const maxChirpLength = 140

	currentChirpLength := len(chirp)
	if currentChirpLength > maxChirpLength {
		return false, "Chirp is too long"
	}

	return true, cleanChirp(chirp)
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
