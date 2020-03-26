package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"firebase.google.com/go/messaging"
	pubsubMiddleware "github.com/fiberweb/pubsub"
	"github.com/gofiber/fiber"

	"worker/config"
	"worker/firebase"
	firestoreEvents "worker/handler/firestore"
	"worker/handler/push"
	"worker/handler/sync"
)

var (
	firebaseApp     *firebase.Firebase
	firestoreClient *firestore.Client
	pubsubClient    *pubsub.Client
	messagingClient *messaging.Client
)

func init() {
	ctx := context.Background()

	var err error
	firebaseApp, err = firebase.New(ctx)
	if err != nil {
		log.Fatalln("Unable to initialize Firebase app:", err)
	}
	firestoreClient, err = firebaseApp.FirestoreClient(ctx)
	if err != nil {
		log.Fatalln("Unable to initialize Firestore client:", err)
	}
	pubsubClient, err = firebaseApp.PubSubClient(ctx)
	if err != nil {
		log.Fatalln("Unable to initialize PubSub client:", err)
	}
	messagingClient, err = firebaseApp.MessagingClient(ctx)
	if err != nil {
		log.Fatalln("Unable to initialize Messaging client:", err)
	}
}

func main() {
	app := fiber.New()

	// all /pubsub/* are to handle PubSub requests
	pubsub := app.Group("/pubsub", pubsubMiddleware.New())
	pubsub.Post("/sync-data", sync.Handler(firestoreClient))
	pubsub.Post("/push-notification", push.Handler(firestoreClient, messagingClient))
	pubsub.Post("/firestore-events", firestoreEvents.Handler(firestoreClient, pubsubClient))
	pubsub.Use(softErrorHandler())

	app.Listen(config.ServicePort)
}
