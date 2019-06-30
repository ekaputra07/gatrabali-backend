package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	goctx "github.com/gorilla/context"
	"github.com/gorilla/mux"
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
		feeds, err := s.db.GetFeeds(context.Background())
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

func (s *server) HandleCategorySummary() http.HandlerFunc {
	// Hardcoded the categories here since we only want to returns these categories (Daerah / Kota) only
	categories := []map[string]interface{}{
		{"id": 2, "title": "Badung"},
		{"id": 3, "title": "Bangli"},
		{"id": 4, "title": "Buleleng"},
		{"id": 5, "title": "Denpasar"},
		{"id": 6, "title": "Gianyar"},
		{"id": 7, "title": "Jembrana"},
		{"id": 8, "title": "Karangasem"},
		{"id": 9, "title": "Klungkung"},
		{"id": 10, "title": "Tabanan"},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var summary []map[string]interface{}

		// loop through categories and get 3 latest entries on that category
		for _, cat := range categories {
			entries, err := s.db.GetCategoryEntries(context.Background(), cat["id"].(int), 0, 3)
			if err != nil {
				continue
			}
			if len(*entries) > 0 {
				cat["entries"] = entries
				summary = append(summary, cat)
			}
		}

		j, err := json.Marshal(summary)
		if err != nil {
			s.SetServerError(w, err.Error())
			return
		}

		if len(summary) > 0 {
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
			entries, err := s.db.GetAllEntries(context.Background(), cur, lim)
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
			entries, err := s.db.GetCategoryEntries(context.Background(), cat, cur, lim)
			if err != nil {
				s.SetServerError(w, err.Error())
				return
			}
			j, err := json.Marshal(entries)
			if err != nil {
				s.SetServerError(w, err.Error())
				return
			}

			// cache 1hr if not empty, if there's user object in request context then only cache for 1 min
			if len(*entries) > 0 {
				if goctx.Get(r, userCtxKey) != nil {
					w = s.SetCacheControl(w, 60)
				} else {
					w = s.SetCacheControl(w, 3600)
				}
			}
			fmt.Fprint(w, string(j))
		}
	}
}

func (s *server) HandleEntry() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		entryID, _ := strconv.Atoi(vars["entryID"])
		entry, err := s.db.GetEntry(context.Background(), entryID)
		if err != nil {
			s.SetServerError(w, err.Error())
			return
		}

		j, err := json.Marshal(entry)
		if err != nil {
			s.SetServerError(w, err.Error())
			return
		}

		w = s.SetCacheControl(w, 3600)
		fmt.Fprint(w, string(j))
	}
}
