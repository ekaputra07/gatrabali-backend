package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
	db     *DB
}

func main() {
	ctx := context.Background()

	// init the Firebase App
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// init the router
	router := mux.NewRouter()

	s := server{router, MakeDB(app)}

	http.Handle("/", s.Routes())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
