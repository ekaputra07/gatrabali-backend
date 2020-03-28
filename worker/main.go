package main

import (
	"context"
	"log"

	pubsubMiddleware "github.com/fiberweb/pubsub"
	"github.com/gofiber/fiber"

	"worker/config"
	"worker/firebase"
	"worker/handler/firestore"
	"worker/handler/push"
	"worker/handler/sync"
)

var firebaseApp *firebase.Firebase

func main() {
	ctx := context.Background()

	// initialize Firebase app
	var err error
	firebaseApp, err = firebase.New(ctx)
	if err != nil {
		log.Fatalln("Unable to initialize Firebase app:", err)
	}

	app := fiber.New()

	// all /pubsub/* are to handle PubSub requests
	pubsub := app.Group("/pubsub", pubsubMiddleware.New())
	pubsub.Post("/sync-data", sync.Handler(ctx, firebaseApp))
	pubsub.Post("/push-notification", push.Handler(ctx, firebaseApp))
	pubsub.Post("/firestore-events", firestore.Handler(ctx, firebaseApp))
	pubsub.Use(softErrorHandler())

	app.Listen(config.ServicePort)
}
