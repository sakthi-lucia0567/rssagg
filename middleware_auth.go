package main

import (
	"fmt"
	"net/http"

	"github.com/sakthi-lucia0567/rssagg/internal/auth"
	internal "github.com/sakthi-lucia0567/rssagg/internal/database"
)

type authHandler func(w http.ResponseWriter, r *http.Request, user internal.User)

func (apiConfig *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
		}
		user, err := apiConfig.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}

		handler(w, r, user)
	}
}
