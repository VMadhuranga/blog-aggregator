package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/VMadhuranga/blog-aggregator/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

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
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		var payload struct {
			Name string
		}
		err := decoder.Decode(&payload)
		if err != nil {
			log.Printf("Error decoding payload: %s", err)
			respondWithError(w, 500, "")
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
		respondWithJson(w, 201, struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Name      string    `json:"lane"`
		}{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Name:      user.Name,
		})
	})

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
