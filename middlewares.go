package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
)

func authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if !strings.HasPrefix(authHeader, "ApiKey ") {
			log.Printf("Error getting api key: %s", errors.New("missing ApiKey prefix"))
			respondWithError(w, 403, "Error getting api key")
			return
		}

		apiKey := strings.Trim(authHeader, "ApiKey ")

		if len(apiKey) == 0 {
			log.Printf("Error getting api key: %s", errors.New("missing ApiKey"))
			respondWithError(w, 403, "Error getting api key")
			return
		}

		ctx := context.WithValue(r.Context(), ctxKey("apiKey"), apiKey)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
