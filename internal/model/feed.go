package model

import "time"

// FeedResponse is the top-level JSON response
type FeedResponse struct {
	Status string    `json:"status"`
	Feed   *FeedMeta `json:"feed"`
	Items  []*Item   `json:"items"`
}

// FeedMeta holds the feed metadata
type FeedMeta struct {
	Title       string     `json:"title"`
	Link        string     `json:"link"`
	Description string     `json:"description"`
	Image       string     `json:"image,omitempty"`
	Updated     *time.Time `json:"updated,omitempty"`
}

// Item represents a single feed entry
type Item struct {
	Title       string     `json:"title"`
	Link        string     `json:"link"`
	Author      string     `json:"author,omitempty"`
	Published   *time.Time `json:"pubDate,omitempty"`
	Content     string     `json:"content,omitempty"`
	Description string     `json:"description,omitempty"`
	Categories  []string   `json:"categories,omitempty"`
	Enclosure   *Enclosure `json:"enclosure,omitempty"`
}

// Enclosure represents media attachments
type Enclosure struct {
	URL    string `json:"url"`
	Type   string `json:"type,omitempty"`
	Length string `json:"length,omitempty"`
}
