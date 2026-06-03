package main

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

func handlerBrowse(s *state, cmd command) error {
	if len(cmd.args) > 1 {
		return fmt.Errorf("browse command requires at most 1 argument")
	}
	limit := 2
	if len(cmd.args) == 1 {
		cmdInt, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("browse command requires an integer argument")
		}
		limit = cmdInt
	}
	posts, err := s.db.GetLatestPosts(context.Background(), int32(limit))
	if err != nil {
		return fmt.Errorf("browse failed: %v", err)
	}
	fmt.Printf("Browsing %d latest posts\n", len(posts))
	for _, post := range posts {
		fmt.Printf("Title: %v\nPublished at: %v\nLink: %v\n%v\n", post.Title.String, post.PublishedAt.Time.Format(time.RFC3339), post.Url.String, post.Description.String)
	}
	return nil
}
