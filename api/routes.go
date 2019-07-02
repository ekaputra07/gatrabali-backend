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
	api.HandleFunc("/entries", s.HandleEntries())
	api.HandleFunc("/entries/{entryID}", s.HandleEntry())
	api.HandleFunc("/categories/summary", s.HandleCategorySummary())
	api.HandleFunc("/kriminal/entries", s.HandleCollectionEntries(constant.Kriminal))
	api.HandleFunc("/baliunited/entries", s.HandleCollectionEntries(constant.BaliUnited))
	api.Use(s.BasicUserCheck)

	return s.router
}
