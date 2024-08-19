package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, statusCode int, errorMessage string) {
	if len(errorMessage) == 0 {
		errorMessage = "Internal server error"
	}

	respondWithJson(w, statusCode, map[string]string{"error": errorMessage})
}

func respondWithJson(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Error encoding response: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if payload != nil {
		w.Write(response)
	}
}
