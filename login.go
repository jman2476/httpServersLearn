package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.Decoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
}
