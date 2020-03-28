package firestore

import (
	"context"
	"errors"
	"net/http"

	"github.com/fiberweb/pubsub"
	"github.com/gofiber/fiber"

	"server/config"
	"server/firebase"
)

// Handler dispatch Firestore events
func Handler(ctx context.Context, fb *firebase.Firebase) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		msg, ok := c.Locals(pubsub.LocalsKey).(*pubsub.Message)
		if !ok {
			c.Next(errors.New("unable to retrieve PubSub message from c.Locals"))
			return
		}

		// check message attributes
		attrs := msg.Message.Attributes
		datatype, ok := attrs["type"]
		if !ok {
			c.Next(errors.New("type is missing from PubSub message attributes"))
			return
		}

		// load clients
		firestoreClient, err := fb.FirestoreClient(ctx)
		if err != nil {
			c.Next(err)
			return
		}
		pubsubClient, err := fb.PubSubClient(ctx, config.GCPProject)
		if err != nil {
			c.Next(err)
			return
		}

		ctx := context.Background() // request ctx

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
