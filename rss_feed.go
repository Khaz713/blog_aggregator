package main

import (
	"blog_aggregator/internal/database"
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, url string) (*RSSFeed, error) {
	httpClient := http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("fetch rss feed: %w", err)
	}
	req.Header.Set("User-Agent", "gator")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch rss feed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch rss feed: bad status: %s", resp.Status)
	}
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fetch rss feed: %w", err)
	}
	var feed RSSFeed
	err = xml.Unmarshal(dat, &feed)
	if err != nil {
		return nil, fmt.Errorf("fetch rss feed: %w", err)
	}
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}
	return &feed, nil
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("scrape feeds: %w", err)
	}

	arg := database.MarkFeedFetchedParams{
		ID:        feed.ID,
		UpdatedAt: time.Now(),
	}

	err = s.db.MarkFeedFetched(context.Background(), arg)
	if err != nil {
		return fmt.Errorf("mark fetched feed: %w", err)
	}

	feedFetched, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	fmt.Printf("Feed %s Fetched, %v posts found\n", feedFetched.Channel.Title, len(feedFetched.Channel.Item))
	for _, item := range feedFetched.Channel.Item {
		pubTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Failed to parse pub date: %s", err)
			pubTime = time.Now()
		}
		arg2 := database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: sql.NullString{
				String: html.UnescapeString(item.Title),
				Valid:  true,
			},
			Url: sql.NullString{
				String: html.UnescapeString(item.Link),
				Valid:  true,
			},
			Description: sql.NullString{
				String: html.UnescapeString(item.Description),
				Valid:  true,
			},
			PublishedAt: sql.NullTime{
				Time:  pubTime,
				Valid: true,
			},
			FeedID: feed.ID,
		}
		_, err = s.db.CreatePost(context.Background(), arg2)
		if err != nil {
			if !strings.Contains(err.Error(), "unique") {
				log.Printf("Failed to create post: %v", err)
			}
		}
	}
	return nil
}
