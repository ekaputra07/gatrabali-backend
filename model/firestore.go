package model

// Category represents category Firestore document
type Category struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

// Feed represents feed Firestore document
type Feed struct {
	ID           int64   `json:"id"`
	UserID       string  `json:"user_id"`
	FeedURL      string  `json:"feed_url"`
	SiteURL      string  `json:"site_url"`
	Title        string  `json:"title"`
	CheckedAt    string  `json:"checked_at"`
	Category     int64   `json:"category"`
	IconID       *int64  `json:"icon_id,omitempty"`
	IconMimeType *string `json:"icon_mime_type,omitempty"`
	IconData     *string `json:"icon_data,omitempty"`
}

// Enclosure is an entry attachment Firestore document
type Enclosure struct {
	URL      string `json:"url"`
	MimeType string `json:"mime_type"`
}

// Entry represent entry Firestore document
type Entry struct {
	ID          int64       `json:"id"`
	UserID      int64       `json:"user_id"`
	FeedID      int64       `json:"feed_id"`
	Hash        string      `json:"hash"`
	Title       string      `json:"title"`
	URL         string      `json:"url"`
	CommentsURL string      `json:"comments_url,omitempty"`
	PublishedAt int64       `json:"published_at"`
	Content     string      `json:"content"`
	Author      string      `json:"author,omitempty"`
	Enclosures  []Enclosure `json:"enclosures,omitempty"`
	Categories  []int64     `json:"categories"`
}
