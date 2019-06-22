package model

// Category represents category Firestore document
type Category struct {
	ID    int64  `json:"id" firestore:"id"`
	Title string `json:"title" firestore:"title"`
}

// CategoryList is a list of Category
type CategoryList []Category

// FeedIcon is the feed icon
type FeedIcon struct {
	ID       int64  `json:"id" firestore:"id"`
	MimeType string `json:"mime_type" firestore:"mime_type"`
	Data     string `json:"data" firestore:"data"`
}

// Feed represents feed Firestore document
type Feed struct {
	ID           int64   `json:"id" firestore:"id"`
	UserID       int64   `json:"user_id" firestore:"user_id"`
	FeedURL      string  `json:"feed_url" firestore:"feed_url"`
	SiteURL      string  `json:"site_url" firestore:"site_url"`
	Title        string  `json:"title" firestore:"title"`
	CheckedAt    string  `json:"checked_at" firestore:"checked_at"`
	Category     int64   `json:"category" firestore:"category"`
	IconID       *int64  `json:"icon_id,omitempty" firestore:"icon_id,omitempty"`
	IconMimeType *string `json:"icon_mime_type,omitempty" firestore:"icon_mime_type,omitempty"`
	IconData     *string `json:"icon_data,omitempty" firestore:"icon_data,omitempty"`
}

// SetIcon sets icon data to Feed object
func (f *Feed) SetIcon(icon *FeedIcon) {
	f.IconID = &icon.ID
	f.IconMimeType = &icon.MimeType
	f.IconData = &icon.Data
}

// Enclosure is an entry attachment Firestore document
type Enclosure struct {
	URL      string `json:"url" firestore:"url"`
	MimeType string `json:"mime_type" firestore:"mime_type"`
}

// Entry represent entry Firestore document
type Entry struct {
	ID          int64        `json:"id" firestore:"id"`
	UserID      int64        `json:"user_id" firestore:"user_id"`
	FeedID      int64        `json:"feed_id" firestore:"feed_id"`
	Title       string       `json:"title" firestore:"title"`
	URL         string       `json:"url" firestore:"url"`
	Content     string       `json:"content" firestore:"content"`
	CommentsURL *string      `json:"comments_url,omitempty" firestore:"comments_url,omitempty"`
	Author      *string      `json:"author,omitempty" firestore:"author,omitempty"`
	Enclosures  *[]Enclosure `json:"enclosures,omitempty" firestore:"enclosures,omitempty"`
	PublishedAt int64        `json:"published_at" firestore:"published_at"`
	Categories  []int64      `json:"categories" firestore:"categories"`
}
