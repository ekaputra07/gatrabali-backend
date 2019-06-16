package model

// MIcon represents Icon on Miniflux
type MIcon struct {
	ID int64 `json:"id"`
}

// MFeed represents Feed on Miniflux
type MFeed struct {
	ID        int64    `json:"id"`
	UserID    string   `json:"user_id"`
	FeedURL   string   `json:"feed_url"`
	SiteURL   string   `json:"site_url"`
	Title     string   `json:"title"`
	CheckedAt string   `json:"checked_at"`
	Category  Category `json:"category"`
	Icon      MIcon    `json:"icon"`
}

// ToFeed transform MFeed to Feed
func (mf *MFeed) ToFeed() (feed Feed) {
	feed = Feed{}
	feed.ID = mf.ID
	feed.UserID = mf.UserID
	feed.SiteURL = mf.SiteURL
	feed.Title = mf.Title
	feed.CheckedAt = mf.CheckedAt
	feed.Category = mf.Category.ID
	return
}

// MEntry is the news item on Miniflux
type MEntry struct {
	ID          int64       `json:"id"`
	UserID      int64       `json:"user_id"`
	FeedID      int64       `json:"feed_id"`
	Hash        string      `json:"hash"`
	Title       string      `json:"title"`
	URL         string      `json:"url"`
	CommentsURL string      `json:"comments_url,omitempty"`
	PublishedAt string      `json:"published_at"`
	Content     string      `json:"content"`
	Author      string      `json:"author,omitempty"`
	Enclosures  []Enclosure `json:"enclosures,omitempty"`
	Feed        MFeed       `json:"feed"`
}

// ToEntry transform MEntry to Entry
func (mf *MEntry) ToEntry() (entry Entry) {
	return
}
