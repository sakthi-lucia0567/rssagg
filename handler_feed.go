package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	internal "github.com/sakthi-lucia0567/rssagg/internal/database"
)

func (apiConfig *apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, user internal.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 403, fmt.Sprintf("Error parsing JSON %v", err))
		return
	}

	userUUID := uuid.New()
	generatedType := pgtype.UUID{Bytes: userUUID, Valid: true}

	feed, err := apiConfig.DB.CreateFeed(r.Context(), internal.CreateFeedParams{
		ID:        generatedType,
		CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Feed not created %v", err))
		return
	}
	respondWithJSON(w, 201, feed)

}
