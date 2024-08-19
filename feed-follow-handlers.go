package main

import (
	"log"
	"net/http"
	"time"

	"github.com/VMadhuranga/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handleCreateFeedFollow(w http.ResponseWriter, r *http.Request) {
	payload, err := decodePayload(r, struct {
		FeedID string `json:"feed_id"`
	}{})

	if err != nil {
		log.Printf("Error decoding payload: %s", err)
		respondWithError(w, 422, "Error decoding payload")
		return
	}

	apiKey := r.Context().Value(ctxKey("apiKey")).(string)
	user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)

	if err != nil {
		log.Printf("Error getting user: %s", err)
		respondWithError(w, 404, "Error getting user")
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    uuid.MustParse(payload.FeedID),
	})

	if err != nil {
		log.Printf("Error creating feed follow: %s", err)
		respondWithError(w, 424, "Error creating feed follow")
		return
	}

	respondWithJson(w, 201, feedFollow)
}

func (apiCfg *apiConfig) handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Context().Value(ctxKey("apiKey")).(string)
	user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)

	if err != nil {
		log.Printf("Error getting user: %s", err)
		respondWithError(w, 404, "Error getting user")
		return
	}

	feedFollowID := r.PathValue("feedFollowID")

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     uuid.MustParse(feedFollowID),
		UserID: user.ID,
	})

	if err != nil {
		log.Printf("Error deleting feed follow: %s", err)
		respondWithError(w, 424, "Error deleting feed follow")
		return
	}

	respondWithJson(w, 204, nil)
}

func (apiCfg *apiConfig) handleGetUserFeedFollows(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Context().Value(ctxKey("apiKey")).(string)
	user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)

	if err != nil {
		log.Printf("Error getting user: %s", err)
		respondWithError(w, 404, "Error getting user")
		return
	}

	userFeedFollows, err := apiCfg.DB.GetUserFeedFollows(r.Context(), user.ID)

	if err != nil {
		log.Printf("Error getting user feed follows: %s", err)
		respondWithError(w, 404, "Error getting user feed follows")
		return
	}

	respondWithJson(w, 200, userFeedFollows)
}
