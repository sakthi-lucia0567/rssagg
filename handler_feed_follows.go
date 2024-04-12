package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	internal "github.com/sakthi-lucia0567/rssagg/internal/database"
)

func (apiConfig *apiConfig) handleCreateFeedFollow(w http.ResponseWriter, r *http.Request, user internal.User) {
	type parameters struct {
		FeedId pgtype.UUID `json:"feed_id"`
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

	feedFollow, err := apiConfig.DB.CreateFeedFollow(r.Context(), internal.CreateFeedFollowParams{
		ID:        generatedType,
		CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		UserID:    user.ID,
		FeedID:    params.FeedId,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Feed not created %v", err))
		return
	}
	respondWithJSON(w, 201, databaseFeedsFollowToFeedFollow(feedFollow))

}

func (apiConfig *apiConfig) handleGetFeedFollow(w http.ResponseWriter, r *http.Request, user internal.User) {

	feedFollow, err := apiConfig.DB.GetFeedFollow(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't Get Feed Follows %v", err))
		return
	}
	respondWithJSON(w, 201, databaseFeedFollowsToFeedFollows(feedFollow))

}

func (apiConfig *apiConfig) handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user internal.User) {
	feedFollowIdStr := chi.URLParam(r, "feedFollowId")
	fmt.Println("feedFollowID", feedFollowIdStr)
	feedFollowId, err := uuid.Parse(feedFollowIdStr)
	fmt.Println("feedFollowID", feedFollowId)
	generatedType := pgtype.UUID{Bytes: feedFollowId, Valid: true}
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't Get FeedFollowId %v", err))
		return
	}

	err = apiConfig.DB.DeleteFeedFollow(r.Context(), internal.DeleteFeedFollowParams{
		ID:     generatedType,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't Delete FeedFollow %v", err))
		return
	}
	respondWithJSON(w, 200, struct{}{})
}
