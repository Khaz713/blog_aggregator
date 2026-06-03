package main

import (
	"errors"
	"fmt"
	"log"
	"time"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("agg requires 1 argument")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("invalid time between requests: %s", err)
	}
	log.Printf("Collecting feeds every %v\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			log.Printf("Failed to scrape feed: %s", err)
		}
	}

	return nil
}
