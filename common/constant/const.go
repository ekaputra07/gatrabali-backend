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
	// Kriminal is collection for crime news entries
	Kriminal = "kriminal"
	// BaleBengong is collection for BaleBengong news entries
	BaleBengong = "balebengong"
	// BaliUnited collection for Bali United news entries
	BaliUnited = "baliunited"
	// EntryResponses collection for responses/comments entries
	EntryResponses = "entry_responses"
)
