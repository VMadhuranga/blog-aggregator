package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/VMadhuranga/blog-aggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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
		log.Fatalf("Error loading environment file: %s", err)
	}

	port := os.Getenv("PORT")
	dbURL := os.Getenv("CONN")
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	apiCfg := apiConfig{
		DB: database.New(db),
	}

	go startFetchingFeeds(10, time.Minute, apiCfg.DB)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Post("/users", apiCfg.handleCreateUser)
	v1Router.Get("/users", authenticate(apiCfg.handleGetUserByApiKey))

	v1Router.Post("/feeds", authenticate(apiCfg.handleCreateFeed))
	v1Router.Get("/feeds", apiCfg.handleGetAllFeeds)

	v1Router.Post("/feed_follows", authenticate(apiCfg.handleCreateFeedFollow))
	v1Router.Get("/feed_follows", authenticate(apiCfg.handleGetUserFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", authenticate(apiCfg.handleDeleteFeedFollow))

	v1Router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		respondWithJson(w, 200, map[string]string{"success": "ok"})
	})

	v1Router.Get("/error", func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, 500, "")
	})

	router.Mount("/v1", v1Router)

	server := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Server running on port: %v", port)
	err = server.ListenAndServe()

	if err != nil {
		log.Fatalf("Error listening on server: %s", err)
	}
}
