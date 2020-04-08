package events

import (
	"context"
	"errors"
	"net/http"

	pubs "github.com/fiberweb/pubsub"
	"github.com/gofiber/fiber"

	"server/common/service"
)

// Handler represents the handler for Firestore events
type Handler struct {
	google *service.Google
}

// New returns Handler instance
func New(g *service.Google) *Handler {
	return &Handler{google: g}
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

		// initialize firestore in case not yet initialize
		if err := h.google.InitFirestore(ctx); err != nil {
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
