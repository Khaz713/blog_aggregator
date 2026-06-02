package main

import (
	"context"
	"errors"
	"log"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("agg expects no arguments")
	}
	url := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}
	log.Print(feed)
	return nil
}
