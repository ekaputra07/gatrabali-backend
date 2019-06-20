package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"api/handler"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handler.Index)

	// Api routes
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/feeds", handler.Feeds).Methods("GET")
	s.HandleFunc("/entries", handler.Entries).Methods("GET")
	s.HandleFunc("/categories/summary", handler.CategorySummary).Methods("GET")

	http.Handle("/", r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
