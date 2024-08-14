package main

import (
	"context"
	"log"
	"net/http"
	"strings"
)

func authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		apiKey := strings.Trim(authHeader, "ApiKey ")
		if !strings.HasPrefix(authHeader, "ApiKey ") || len(apiKey) == 0 {
			log.Println("Error getting api key")
			respondWithError(w, 403, "Invalid api key")
			return
		}
		ctx := context.WithValue(r.Context(), ctxKey("apiKey"), apiKey)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
