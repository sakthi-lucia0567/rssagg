package main

import (
	"time"

	"github.com/google/uuid"
	internal "github.com/sakthi-lucia0567/rssagg/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func databaseUserToUser(dbUser internal.User) User {
	return User{
		ID:        dbUser.ID.Bytes,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
		Name:      dbUser.Name,
		ApiKey:    dbUser.ApiKey,
	}
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserId    uuid.UUID `json:"user_id"`
}

func databaseFeedToFeed(dbFeed internal.Feed) Feed {
	return Feed{
		ID:        dbFeed.ID.Bytes,
		CreatedAt: dbFeed.CreatedAt.Time,
		UpdatedAt: dbFeed.UpdatedAt.Time,
		Name:      dbFeed.Name,
		Url:       dbFeed.Url,
		UserId:    dbFeed.UserID.Bytes,
	}
}

func databaseFeedsToFeeds(dbFeeds []internal.Feed) []Feed {
	feeds := []Feed{}
	for _, dbFeed := range dbFeeds {
		feeds = append(feeds, databaseFeedToFeed(dbFeed))
	}
	return feeds
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func databaseFeedsFollowToFeedFollow(dbFeedFollow internal.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        dbFeedFollow.ID.Bytes,
		CreatedAt: dbFeedFollow.CreatedAt.Time,
		UpdatedAt: dbFeedFollow.UpdatedAt.Time,
		UserID:    dbFeedFollow.UserID.Bytes,
		FeedID:    dbFeedFollow.FeedID.Bytes,
	}
}

func databaseFeedFollowsToFeedFollows(dbFeedFollows []internal.FeedFollow) []FeedFollow {
	feedFollows := []FeedFollow{}
	for _, dbFeedFollow := range dbFeedFollows {
		feedFollows = append(feedFollows, databaseFeedsFollowToFeedFollow(dbFeedFollow))
	}
	return feedFollows
}
