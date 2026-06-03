package main

import (
	"blog_aggregator/internal/database"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return errors.New("addfeed expects 2 arguments")
	}

	arg := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), arg)
	if err != nil {
		return fmt.Errorf("create feed: %w", err)
	}
	log.Printf("Feed created; Name: %s, URL: %s", feed.Name, feed.Url)

	arg2 := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), arg2)
	if err != nil {
		return fmt.Errorf("create feed follow: %w", err)
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("feeds expects no arguments")
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("get feeds: %w", err)
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("get user: %w", err)
		}
		log.Printf("Feed found; Name: %s, URL: %s, Created By: %s\n", feed.Name, feed.Url, user.Name)
	}
	return nil
}
