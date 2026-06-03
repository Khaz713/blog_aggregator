package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"log"
	"time"

	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("follow expects 1 argument")
	}
	url := cmd.args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}

	arg := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), arg)
	if err != nil {
		return fmt.Errorf("error creating follow: %w", err)
	}

	log.Printf("User (%s) followed feed: %s", user.Name, feed.Name)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("following expects no arguments")
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error fetching follows: %w", err)
	}

	log.Printf("User (%s) follows %d feeds:\n", user.Name, len(follows))
	for _, follow := range follows {
		feed, err := s.db.GetFeedById(context.Background(), follow.FeedID)
		if err != nil {
			return fmt.Errorf("error fetching feed: %w", err)
		}
		log.Printf("* %s\n", feed.Name)
	}
	return nil
}

func handlerUnFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("unfollow expects 1 argument")
	}
	url := cmd.args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}
	feedId := feed.ID
	arg := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feedId,
	}
	err = s.db.DeleteFeedFollow(context.Background(), arg)
	if err != nil {
		return fmt.Errorf("error deleting follow: %w", err)
	}
	log.Printf("User (%s) unfollowed feed: %s", user.Name, feed.Name)
	return nil
}
