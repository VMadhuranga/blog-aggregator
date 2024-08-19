package main

import (
	"log"
	"net/http"
	"time"

	"github.com/VMadhuranga/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	payload, err := decodePayload(r, struct {
		Name string
	}{})

	if err != nil {
		log.Printf("Error decoding payload: %s", err)
		respondWithError(w, 422, "Error decoding payload")
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      payload.Name,
	})

	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(w, 424, "Error creating user")
		return
	}

	respondWithJson(w, 201, user)
}

func (apiCfg *apiConfig) handleGetUserByApiKey(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Context().Value(ctxKey("apiKey")).(string)
	user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)

	if err != nil {
		log.Printf("Error getting user: %s", err)
		respondWithError(w, 404, "Error getting user")
		return
	}

	respondWithJson(w, 200, user)
}
