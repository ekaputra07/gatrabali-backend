package model

import (
	"errors"
	"strings"
	"time"
)

// MFeed represents Feed on Miniflux
type MFeed struct {
	ID        int64    `json:"id"`
	UserID    int64    `json:"user_id"`
	FeedURL   string   `json:"feed_url"`
	SiteURL   string   `json:"site_url"`
	Title     string   `json:"title"`
	CheckedAt string   `json:"checked_at"`
	Category  Category `json:"category"`
}

// ToFeed transform MFeed into Feed
func (mf *MFeed) ToFeed() (feed Feed) {
	feed = Feed{
		ID:        mf.ID,
		UserID:    mf.UserID,
		FeedURL:   mf.FeedURL,
		SiteURL:   mf.SiteURL,
		Title:     mf.Title,
		CheckedAt: mf.CheckedAt,
		Category:  mf.Category.ID,
	}
	return
}

// MEntry is the news item on Miniflux
type MEntry struct {
	ID          int64        `json:"id"`
	UserID      int64        `json:"user_id"`
	FeedID      int64        `json:"feed_id"`
	Title       string       `json:"title"`
	URL         string       `json:"url"`
	Content     string       `json:"content"`
	CommentsURL *string      `json:"comments_url,omitempty"`
	Author      *string      `json:"author,omitempty"`
	Enclosures  *[]Enclosure `json:"enclosures,omitempty"`
	PublishedAt string       `json:"published_at"`
	Feed        MFeed        `json:"feed"`
}

// ToEntry transform MEntry into Entry
func (me *MEntry) ToEntry() (*Entry, error) {

	// return error if enclosure nil
	if me.Enclosures == nil || (me.Enclosures != nil && len(*me.Enclosures) == 0) {
		return nil, errors.New("Entry doesn't have enclosure")
	}
	// Return error if enclosure is not an image
	if !strings.Contains((*me.Enclosures)[0].MimeType, "image") {
		return nil, errors.New("Enclosure is not an image")
	}

	entry := Entry{
		ID:          me.ID,
		UserID:      me.UserID,
		FeedID:      me.FeedID,
		Title:       me.Title,
		URL:         me.URL,
		Content:     me.Content,
		CommentsURL: me.CommentsURL,
		Author:      me.Author,
		Enclosures:  me.Enclosures,
	}
	// transform published at to unix timestamp
	t, _ := time.Parse(time.RFC3339, me.PublishedAt)
	entry.PublishedAt = t.Unix() * 1000 // in millisecs

	// set category
	entry.CategoryID = me.Feed.Category.ID
	entry.Categories = []int64{me.Feed.Category.ID} // to support legacy app.

	return &entry, nil
}
