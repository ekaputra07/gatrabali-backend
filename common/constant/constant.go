package constant

const (
	// TypeCategory news category
	TypeCategory = "CATEGORY"
	// TypeFeed feed source
	TypeFeed = "FEED"
	// TypeEntry news entry
	TypeEntry = "ENTRY"

	// OpWrite is write operation on Firestore
	OpWrite = "WRITE"
	// OpDelete is delete operation on Firestore
	OpDelete = "DELETE"

	// Categories is collection for categories
	Categories = "categories"
	// Feeds is collection for feed sources
	Feeds = "feeds"
	// Entries is collection for news entries
	Entries = "entries"
)
