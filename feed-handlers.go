package main

import (
	"log"
	"net/http"
	"time"

	"github.com/VMadhuranga/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request) {
	payload, err := decodePayload(r, struct {
		Name string
		Url  string
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

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      payload.Name,
		Url:       payload.Url,
		UserID:    user.ID,
	})

	if err != nil {
		log.Printf("Error creating feed: %s", err)
		respondWithError(w, 424, "Error creating feed")
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		log.Printf("Error creating feed follow: %s", err)
		respondWithError(w, 424, "Error creating feed follow")
		return
	}

	respondWithJson(w, 201, struct {
		Feed       database.Feed       `json:"feed"`
		FeedFollow database.FeedFollow `json:"feed_follow"`
	}{
		Feed: database.Feed{
			ID:        feed.ID,
			CreatedAt: feed.CreatedAt,
			UpdatedAt: feed.UpdatedAt,
			Name:      feed.Name,
			Url:       feed.Url,
			UserID:    user.ID,
		},
		FeedFollow: database.FeedFollow{
			ID:        feedFollow.ID,
			CreatedAt: feedFollow.CreatedAt,
			UpdatedAt: feedFollow.UpdatedAt,
			UserID:    feedFollow.UserID,
			FeedID:    feedFollow.FeedID,
		},
	})
}

func (apiCfg *apiConfig) handleGetAllFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetAllFeeds(r.Context())

	if err != nil {
		log.Printf("Error getting all feeds: %s", err)
		respondWithError(w, 404, "Error getting all feeds")
		return
	}

	respondWithJson(w, 200, feeds)
}
