package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
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
