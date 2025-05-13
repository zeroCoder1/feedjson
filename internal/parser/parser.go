package parser

import (
	"context"
	"errors"

	"github.com/mmcdole/gofeed"

	"github.com/zeroCoder1/feedjson/internal/model"
)

var (
	// ErrInvalidURL is returned when rss_url query param is missing
	ErrInvalidURL = errors.New("invalid rss_url parameter")
)

// FetchFeed retrieves, parses, and returns a FeedResponse
func FetchFeed(ctx context.Context, rssURL string, count int) (*model.FeedResponse, error) {
	if rssURL == "" {
		return nil, ErrInvalidURL
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(rssURL, ctx)
	if err != nil {
		return nil, err
	}

	items := make([]*model.Item, 0, len(feed.Items))
	for i, it := range feed.Items {
		if count > 0 && i >= count {
			break
		}
		item := &model.Item{
			Title:       it.Title,
			Link:        it.Link,
			Author:      "",
			Published:   it.PublishedParsed,
			Content:     it.Content,
			Description: it.Description,
			Categories:  it.Categories,
		}
		if it.Author != nil {
			item.Author = it.Author.Name
		}
		if len(it.Enclosures) > 0 {
			enc := it.Enclosures[0]
			item.Enclosure = &model.Enclosure{
				URL:    enc.URL,
				Type:   enc.Type,
				Length: enc.Length,
			}
		}
		items = append(items, item)
	}

	resp := &model.FeedResponse{
		Status: "ok",
		Feed: &model.FeedMeta{
			Title:       feed.Title,
			Link:        feed.Link,
			Description: feed.Description,
			Updated:     feed.UpdatedParsed,
		},
		Items: items,
	}
	if feed.Image != nil {
		resp.Feed.Image = feed.Image.URL
	}

	return resp, nil
}
