package firestore

import (
	"context"
	"errors"
	"net/http"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	pubs "github.com/fiberweb/pubsub"
	"github.com/gofiber/fiber"

	"server/config"
	"server/firebase"
)

// Handler represents the handler for Firestore events
type Handler struct {
	fb *firebase.Firebase

	firestore *firestore.Client // lazily set
	pubsub    *pubsub.Client    // lazily set
}

// New returns Handler instance
func New(fb *firebase.Firebase) *Handler {
	return &Handler{fb: fb}
}

// Handle handles the request
func (h *Handler) Handle() func(*fiber.Ctx) {

	return func(c *fiber.Ctx) {
		msg, ok := c.Locals(pubs.LocalsKey).(*pubs.Message)
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

		ctx := context.Background()

		// load clients
		var err error
		h.firestore, err = h.fb.FirestoreClient(ctx)
		if err != nil {
			c.Next(err)
			return
		}
		h.pubsub, err = h.fb.PubSubClient(ctx, config.GCPProject)
		if err != nil {
			c.Next(err)
			return
		}

		switch datatype.(string) {
		case "entries":
			if err := h.notifySubscribers(ctx, msg.Message.Data); err != nil {
				c.Next(err)
				return
			}
		case "responses":
			if err := h.aggregateResponses(ctx, msg.Message.Data); err != nil {
				c.Next(err)
				return
			}
		}
		c.SendStatus(http.StatusOK)
	}
}
