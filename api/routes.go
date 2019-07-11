package main

import (
	"github.com/apps4bali/gatrabali-backend/common/constant"
	"github.com/gorilla/mux"
)

func (s *server) Routes() *mux.Router {
	// index route
	s.router.HandleFunc("/", s.HandleIndex())

	// Api routes
	api := s.router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/feeds", s.HandleFeeds())
	api.HandleFunc("/entries", s.HandleEntries(constant.Entries))
	api.HandleFunc("/entries/{entryID}", s.HandleEntry(constant.Entries))
	api.HandleFunc("/categories/summary", s.HandleCategorySummary(constant.Entries, 3))
	api.HandleFunc("/kriminal/entries", s.HandleEntries(constant.Kriminal))
	api.HandleFunc("/kriminal/entries/{entryID}", s.HandleEntry(constant.Kriminal))
	api.HandleFunc("/baliunited/entries", s.HandleEntries(constant.BaliUnited))
	api.HandleFunc("/baliunited/entries/{entryID}", s.HandleEntry(constant.BaliUnited))
	api.HandleFunc("/balebengong/entries", s.HandleEntries(constant.BaleBengong))
	api.HandleFunc("/balebengong/entries/{entryID}", s.HandleEntry(constant.BaleBengong))
	api.HandleFunc("/balebengong/summary", s.HandleCategorySummary(constant.BaleBengong, 3))

	api.Use(s.BasicUserCheck)

	return s.router
}
