package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, errMSG error) {
	logMsg := fmt.Sprintf("Sending: %d-%s", code, msg)
	if errMSG != nil {
		logMsg = fmt.Sprintf("Encountered: %s\n		", errMSG) + logMsg
	}
	log.Println(logMsg)

	errResp := struct {
		Error string `json:"error"`
	}{
		Error: msg,
	}

	respondWithJSON(w, code, errResp)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error Marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}
