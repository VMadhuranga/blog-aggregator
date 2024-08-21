package main

import (
	"log"
	"net/http"

	"github.com/VMadhuranga/blog-aggregator/internal/database"
)

func (apiCfg *apiConfig) handleGetPostsByUser(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Context().Value(ctxKey("apiKey")).(string)
	user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)

	if err != nil {
		log.Printf("Error getting user: %s", err)
		respondWithError(w, 404, "Error getting user")
		return
	}

	postsByUser, err := apiCfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  10,
	})

	if err != nil {
		log.Printf("Error getting posts by user: %s", err)
		respondWithError(w, 404, "Error getting posts by user")
		return
	}

	respondWithJson(w, 200, postsByUser)
}
