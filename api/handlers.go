package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (s *server) SetCacheControl(w http.ResponseWriter, maxAge int) http.ResponseWriter {
	w.Header().Add("Cache-Control", fmt.Sprintf("public, max-age=%v, s-maxage=%v", maxAge, maxAge))
	return w
}

func (s *server) SetServerError(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusInternalServerError)
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
			s.SetServerError(w, err.Error())
			return
		}
		j, err := json.Marshal(feeds)
		if err != nil {
			s.SetServerError(w, err.Error())
			return
		}

		// cache if not empty
		if len(*feeds) > 0 {
			w = s.SetCacheControl(w, 3600)
		}
		fmt.Fprint(w, string(j))
	}
}

func (s *server) HandleEntries() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		cat, _ := strconv.Atoi(query.Get("categoryId"))
		cur, _ := strconv.Atoi(query.Get("cursor"))
		lim, _ := strconv.Atoi(query.Get("limit"))

		if cat == 0 {
			// Returns latest entries
			entries, err := s.db.GetAllEntries(cur, lim)
			if err != nil {
				s.SetServerError(w, err.Error())
				return
			}
			j, err := json.Marshal(entries)
			if err != nil {
				s.SetServerError(w, err.Error())
				return
			}

			// cache if not empty
			if len(*entries) > 0 {
				w = s.SetCacheControl(w, 3600)
			}
			fmt.Fprint(w, string(j))

		} else {
			// Returns latest entries in category
			entries, err := s.db.GetCategoryEntries(cat, cur, lim)
			if err != nil {
				s.SetServerError(w, err.Error())
				return
			}
			j, err := json.Marshal(entries)
			if err != nil {
				s.SetServerError(w, err.Error())
				return
			}

			// cache if not empty
			if len(*entries) > 0 {
				w = s.SetCacheControl(w, 3600)
			}
			fmt.Fprint(w, string(j))
		}
	}
}
