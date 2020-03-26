package firestore

import (
	"context"
	"errors"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	pubs "github.com/fiberweb/pubsub"
	"github.com/gofiber/fiber"
)

// Handler dispatch Firestore events
func Handler(firestoreClient *firestore.Client, pubsubClient *pubsub.Client) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		ctx := context.Background()

		msg := c.Locals(pubs.LocalsKey).(*pubs.Message)
		log.Printf("Firestore event received with payload: %v\n", msg)

		// check message attributes
		attrs := msg.Message.Attributes
		datatype, ok := attrs["type"]
		if !ok {
			c.Next(errors.New("type is missing from PubSub message attributes"))
			return
		}

		switch datatype.(string) {
		case "entries":
			if err := notifySubscriber(ctx, firestoreClient, pubsubClient, msg.Message.Data); err != nil {
				c.Next(err)
				return
			}
		case "responses":
			if err := aggregateEntryResponse(ctx, firestoreClient, msg.Message.Data); err != nil {
				c.Next(err)
				return
			}
		}

		c.SendStatus(http.StatusOK)
	}
}
