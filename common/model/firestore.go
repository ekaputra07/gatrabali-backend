package model

import "encoding/json"

// Category represents category Firestore document
type Category struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

// ToMap converts this struct into a map with lowercase keys (store this format on Firestore)
func (c *Category) ToMap() (*map[string]interface{}, error) {
	j, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(j, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// CategoryList is a list of Category
type CategoryList []Category

// FeedIcon is the feed icon
type FeedIcon struct {
	ID       int64  `json:"id"`
	MimeType string `json:"mime_type"`
	Data     string `json:"data"`
}

// Feed represents feed Firestore document
type Feed struct {
	ID           int64   `json:"id"`
	UserID       int64   `json:"user_id"`
	FeedURL      string  `json:"feed_url"`
	SiteURL      string  `json:"site_url"`
	Title        string  `json:"title"`
	CheckedAt    string  `json:"checked_at"`
	Category     int64   `json:"category"`
	IconID       *int64  `json:"icon_id,omitempty"`
	IconMimeType *string `json:"icon_mime_type,omitempty"`
	IconData     *string `json:"icon_data,omitempty"`
}

// SetIcon sets icon data to Feed object
func (f *Feed) SetIcon(icon *FeedIcon) {
	f.IconID = &icon.ID
	f.IconMimeType = &icon.MimeType
	f.IconData = &icon.Data
}

// ToMap converts this struct into a map with lowercase keys (store this format on Firestore)
func (f *Feed) ToMap() (*map[string]interface{}, error) {
	j, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(j, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// Enclosure is an entry attachment Firestore document
type Enclosure struct {
	URL      string `json:"url"`
	MimeType string `json:"mime_type"`
}

// Entry represent entry Firestore document
type Entry struct {
	ID          int64        `json:"id"`
	UserID      int64        `json:"user_id"`
	FeedID      int64        `json:"feed_id"`
	Hash        string       `json:"hash"`
	Title       string       `json:"title"`
	URL         string       `json:"url"`
	Content     string       `json:"content"`
	CommentsURL *string      `json:"comments_url,omitempty"`
	Author      *string      `json:"author,omitempty"`
	Enclosures  *[]Enclosure `json:"enclosures,omitempty"`
	PublishedAt int64        `json:"published_at"`
	Categories  []int64      `json:"categories"`
}

// ToMap converts this struct into a map with lowercase keys (store this format on Firestore)
func (e *Entry) ToMap() (*map[string]interface{}, error) {
	j, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(j, &m); err != nil {
		return nil, err
	}
	return &m, nil
}
