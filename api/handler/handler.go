package handler

import (
	"fmt"
	"net/http"
)

// Index serve request to /
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Gatra Bali API v1.0.0")
}

// Feeds serve request to /api/v1/feeds
func Feeds(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Feeds API")
}

// Entries serve request to /api/v1/entries
func Entries(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Entries API")
}

// CategorySummary serve request to /api/v1/category_summary
func CategorySummary(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "CategorySummary API")
}
