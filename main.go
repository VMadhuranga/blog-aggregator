package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading environment file: %s", err)
		return
	}
	port := os.Getenv("PORT")
	serveMux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

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
