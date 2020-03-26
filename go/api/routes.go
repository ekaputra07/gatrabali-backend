package main

import (
	"github.com/apps4bali/gatrabali-backend/go/common"
	"github.com/gorilla/mux"
)

func (s *server) Routes() *mux.Router {
	// index route
	s.router.HandleFunc("/", s.HandleIndex())

	// Api V1 routes:
	// - entries returned by this API version is ordered by entry ID descending
	api := s.router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/feeds", s.HandleFeeds())
	api.HandleFunc("/entries", s.HandleEntries(common.Entries, "id"))
	api.HandleFunc("/entries/{entryID}", s.HandleEntry(common.Entries))
	api.HandleFunc("/categories/summary", s.HandleCategorySummary(common.Entries, "id", 3))
	api.HandleFunc("/kriminal/entries", s.HandleEntries(common.Kriminal, "id"))
	api.HandleFunc("/kriminal/entries/{entryID}", s.HandleEntry(common.Kriminal))
	api.HandleFunc("/baliunited/entries", s.HandleEntries(common.BaliUnited, "id"))
	api.HandleFunc("/baliunited/entries/{entryID}", s.HandleEntry(common.BaliUnited))
	// api.Use(s.BasicUserCheck)

	// Api V2 routes:
	// - entries returned by this API version is ordered by entry 'published_at' descending
	// - added support for BaleBengong entries
	apiV2 := s.router.PathPrefix("/api/v2").Subrouter()
	apiV2.HandleFunc("/feeds", s.HandleFeeds())
	apiV2.HandleFunc("/entries", s.HandleEntries(common.Entries, "published_at"))
	apiV2.HandleFunc("/entries/{entryID}", s.HandleEntry(common.Entries))
	apiV2.HandleFunc("/categories/summary", s.HandleCategorySummary(common.Entries, "published_at", 3))
	apiV2.HandleFunc("/kriminal/entries", s.HandleEntries(common.Kriminal, "published_at"))
	apiV2.HandleFunc("/kriminal/entries/{entryID}", s.HandleEntry(common.Kriminal))
	apiV2.HandleFunc("/baliunited/entries", s.HandleEntries(common.BaliUnited, "published_at"))
	apiV2.HandleFunc("/baliunited/entries/{entryID}", s.HandleEntry(common.BaliUnited))
	apiV2.HandleFunc("/balebengong/entries", s.HandleEntries(common.BaleBengong, "published_at"))
	apiV2.HandleFunc("/balebengong/entries/{entryID}", s.HandleEntry(common.BaleBengong))
	apiV2.HandleFunc("/balebengong/summary", s.HandleCategorySummary(common.BaleBengong, "published_at", 3))

	return s.router
}
