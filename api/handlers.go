package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"server/common/constant"
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
		fmt.Fprint(w, "Gatra Bali Backend: https://server")
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

func (s *server) HandleEntry(collection string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		entryID, _ := strconv.Atoi(vars["entryID"])
		entry, err := s.db.GetCollectionEntry(context.Background(), collection, entryID)
		if err != nil {
			s.SetServerError(w, err.Error())
			return
		}

		j, err := json.Marshal(entry)
		if err != nil {
			s.SetServerError(w, err.Error())
			return
		}

		w = s.SetCacheControl(w, 86400) // cache 24hr
		fmt.Fprint(w, string(j))
	}
}

func (s *server) HandleCategorySummary(collection, orderBy string, limit int) http.HandlerFunc {
	var categories []map[string]interface{}

	// Hardcoded the categories here since we only want to returns these categories for specified collection
	if collection == constant.Entries {
		categories = []map[string]interface{}{
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
	} else if collection == constant.BaleBengong {
		categories = []map[string]interface{}{
			{"id": 13, "title": "Opini"},
			{"id": 14, "title": "Teknologi"},
			{"id": 15, "title": "Lingkungan"},
			{"id": 16, "title": "Sosok"},
			{"id": 17, "title": "Budaya"},
			{"id": 18, "title": "Sosial"},
			{"id": 19, "title": "Agenda"},
			{"id": 20, "title": "Travel"},
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var summary []map[string]interface{}

		// loop through categories and get 3 latest entries on that category
		for _, cat := range categories {
			entries, err := s.db.GetCollectionCategoryEntries(context.Background(), collection, cat["id"].(int), orderBy, 0, limit)
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

func (s *server) HandleEntries(collection, orderBy string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		cat, _ := strconv.Atoi(query.Get("categoryId"))
		cur, _ := strconv.Atoi(query.Get("cursor"))
		lim, _ := strconv.Atoi(query.Get("limit"))

		if cat == 0 {
			// Returns latest entries
			entries, err := s.db.GetCollectionEntries(context.Background(), collection, orderBy, cur, lim)
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
			entries, err := s.db.GetCollectionCategoryEntries(context.Background(), collection, cat, orderBy, cur, lim)
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
				w = s.SetCacheControl(w, 3600)
			}
			fmt.Fprint(w, string(j))
		}
	}
}
