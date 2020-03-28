package main

import (
	"server/common/constant"
	"github.com/gorilla/mux"
)

func (s *server) Routes() *mux.Router {
	// index route
	s.router.HandleFunc("/", s.HandleIndex())

	// Api V1 routes:
	// - entries returned by this API version is ordered by entry ID descending
	api := s.router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/feeds", s.HandleFeeds())
	api.HandleFunc("/entries", s.HandleEntries(constant.Entries, "id"))
	api.HandleFunc("/entries/{entryID}", s.HandleEntry(constant.Entries))
	api.HandleFunc("/categories/summary", s.HandleCategorySummary(constant.Entries, "id", 3))
	api.HandleFunc("/kriminal/entries", s.HandleEntries(constant.Kriminal, "id"))
	api.HandleFunc("/kriminal/entries/{entryID}", s.HandleEntry(constant.Kriminal))
	api.HandleFunc("/baliunited/entries", s.HandleEntries(constant.BaliUnited, "id"))
	api.HandleFunc("/baliunited/entries/{entryID}", s.HandleEntry(constant.BaliUnited))
	// api.Use(s.BasicUserCheck)

	// Api V2 routes:
	// - entries returned by this API version is ordered by entry 'published_at' descending
	// - added support for BaleBengong entries
	apiV2 := s.router.PathPrefix("/api/v2").Subrouter()
	apiV2.HandleFunc("/feeds", s.HandleFeeds())
	apiV2.HandleFunc("/entries", s.HandleEntries(constant.Entries, "published_at"))
	apiV2.HandleFunc("/entries/{entryID}", s.HandleEntry(constant.Entries))
	apiV2.HandleFunc("/categories/summary", s.HandleCategorySummary(constant.Entries, "published_at", 3))
	apiV2.HandleFunc("/kriminal/entries", s.HandleEntries(constant.Kriminal, "published_at"))
	apiV2.HandleFunc("/kriminal/entries/{entryID}", s.HandleEntry(constant.Kriminal))
	apiV2.HandleFunc("/baliunited/entries", s.HandleEntries(constant.BaliUnited, "published_at"))
	apiV2.HandleFunc("/baliunited/entries/{entryID}", s.HandleEntry(constant.BaliUnited))
	apiV2.HandleFunc("/balebengong/entries", s.HandleEntries(constant.BaleBengong, "published_at"))
	apiV2.HandleFunc("/balebengong/entries/{entryID}", s.HandleEntry(constant.BaleBengong))
	apiV2.HandleFunc("/balebengong/summary", s.HandleCategorySummary(constant.BaleBengong, "published_at", 3))

	return s.router
}
