package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *server) SetCacheControl(w http.ResponseWriter, maxAge int) http.ResponseWriter {
	w.Header().Add("Cache-Control", fmt.Sprintf("public, max-age=%v, s-maxage=%v", maxAge, maxAge))
	return w
}

func (s *server) HandleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Gatra Bali Backend: https://github.com/apps4bali/gatrabali-backend")
	}
}

func (s *server) HandleFeeds() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		feeds, err := s.db.GetFeeds()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(feeds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(s.SetCacheControl(w, 3600), string(j))
	}
}
