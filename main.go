package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/VMadhuranga/blog-aggregator/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type ctxKey string
type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading environment file: %s", err)
		return
	}
	port := os.Getenv("PORT")
	dbURL := os.Getenv("CONN")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error opening database: %s", err)
		return
	}
	dbQueries := database.New(db)
	cnfg := apiConfig{
		DB: dbQueries,
	}
	ctx := context.Background()
	serveMux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	// create user handler
	serveMux.HandleFunc("POST /v1/users", func(w http.ResponseWriter, r *http.Request) {
		payload, err := decodePayload(r, struct {
			Name string
		}{})
		if err != nil {
			log.Printf("Error decoding payload: %s", err)
			respondWithError(w, 500, "")
			return
		}
		user, err := cnfg.DB.CreateUser(ctx, database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      payload.Name,
		})
		if err != nil {
			log.Printf("Error creating user: %s", err)
			respondWithError(w, 500, "")
			return
		}
		user.ApiKey = "" // omit api key from response
		respondWithJson(w, 201, user)
	})

	// get user by api key handler
	serveMux.HandleFunc("GET /v1/users", authenticate(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Context().Value(ctxKey("apiKey")).(string)
		user, err := cnfg.DB.GetUserByApiKey(ctx, apiKey)
		if err != nil {
			log.Printf("Error getting user: %s", err)
			respondWithError(w, 500, "")
			return
		}
		respondWithJson(w, 200, user)
	}))

	// create feed handler
	serveMux.HandleFunc("POST /v1/feeds", authenticate(func(w http.ResponseWriter, r *http.Request) {
		payload, err := decodePayload(r, struct {
			Name string
			Url  string
		}{})
		if err != nil {
			log.Printf("Error decoding payload: %s", err)
			respondWithError(w, 500, "")
			return
		}
		apiKey := r.Context().Value(ctxKey("apiKey")).(string)
		user, err := cnfg.DB.GetUserByApiKey(ctx, apiKey)
		if err != nil {
			log.Printf("Error getting user: %s", err)
			respondWithError(w, 500, "")
			return
		}
		feed, err := cnfg.DB.CreateFeed(ctx, database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      payload.Name,
			Url:       payload.Url,
			UserID:    user.ID,
		})
		if err != nil {
			log.Printf("Error creating feed: %s", err)
			respondWithError(w, 500, "")
			return
		}
		feedFollow, err := cnfg.DB.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		})
		if err != nil {
			log.Printf("Error creating feed follow: %s", err)
			respondWithError(w, 500, "")
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
	}))

	// get all feeds handler
	serveMux.HandleFunc("GET /v1/feeds", func(w http.ResponseWriter, r *http.Request) {
		feeds, err := cnfg.DB.GetAllFeeds(ctx)
		if err != nil {
			log.Printf("Error getting all feeds: %s", err)
			respondWithError(w, 500, "")
			return
		}
		respondWithJson(w, 200, feeds)
	})

	// create feed follow handler
	serveMux.HandleFunc("POST /v1/feed_follows", authenticate(func(w http.ResponseWriter, r *http.Request) {
		payload, err := decodePayload(r, struct {
			FeedID string `json:"feed_id"`
		}{})
		if err != nil {
			log.Printf("Error decoding payload: %s", err)
			respondWithError(w, 500, "")
			return
		}
		apiKey := r.Context().Value(ctxKey("apiKey")).(string)
		user, err := cnfg.DB.GetUserByApiKey(ctx, apiKey)
		if err != nil {
			log.Printf("Error getting user: %s", err)
			respondWithError(w, 500, "")
			return
		}
		feedFollow, err := cnfg.DB.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    uuid.MustParse(payload.FeedID),
		})
		if err != nil {
			log.Printf("Error creating feed follow: %s", err)
			respondWithError(w, 500, "")
			return
		}
		respondWithJson(w, 201, feedFollow)
	}))

	// delete feed follow handler
	serveMux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", func(w http.ResponseWriter, r *http.Request) {
		feedFollowID := r.PathValue("feedFollowID")
		err := cnfg.DB.DeleteFeedFollow(ctx, uuid.MustParse(feedFollowID))
		if err != nil {
			log.Panicf("Error deleting feed flow: %s", err)
			respondWithError(w, 500, "")
			return
		}
		respondWithJson(w, 204, nil)
	})

	// get user feed follows handler
	serveMux.HandleFunc("GET /v1/feed_follows", authenticate(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Context().Value(ctxKey("apiKey")).(string)
		user, err := cnfg.DB.GetUserByApiKey(ctx, apiKey)
		if err != nil {
			log.Printf("Error getting user: %s", err)
			respondWithError(w, 500, "")
			return
		}
		userFeedFollows, err := cnfg.DB.GetUserFeedFollows(ctx, user.ID)
		if err != nil {
			log.Printf("Error getting user feed follows: %s", err)
			respondWithError(w, 500, "")
			return
		}
		respondWithJson(w, 200, userFeedFollows)
	}))

	// test respondWithJson function
	serveMux.HandleFunc("GET /v1/healthz", func(w http.ResponseWriter, r *http.Request) {
		respondWithJson(w, 200, map[string]string{"success": "ok"})
	})

	// test respondWithError function
	serveMux.HandleFunc("GET /v1/error", func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, 500, "")
	})

	err = server.ListenAndServe()
	if err != nil {
		log.Printf("Error listening on server: %s", err)
		return
	}
}
