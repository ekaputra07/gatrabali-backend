package firestore

import (
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"github.com/gofiber/fiber"
)

// Handler dispatch Firestore events
func Handler(firestoreClient *firestore.Client, pubsubClient *pubsub.Client) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {

	}
}
