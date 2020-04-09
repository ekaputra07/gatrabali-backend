package sync

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fiberweb/pubsub"
	"github.com/gofiber/fiber"

	"server/common/constant"
	"server/common/service"
	"server/common/types"
)

// Handler represents the data syncer from Miniflux to Firestore
type Handler struct {
	google *service.Google
}

// New returns an instance of Handler
func New(google *service.Google) *Handler {
	return &Handler{google}
}

// Handle handles the request
func (h *Handler) Handle() func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		msg, ok := c.Locals(pubsub.LocalsKey).(*pubsub.Message)
		if !ok {
			c.Next(errors.New("unable to retrieve PubSub message from c.Locals"))
			return
		}

		var payload *types.SyncPayload
		if err := json.Unmarshal(msg.Message.Data, &payload); err != nil {
			c.Next(err)
			return
		}
		if payload.ID == nil || payload.Type == nil || payload.Op == nil {
			c.Next(errors.New("Invalid message payload: missing id, type or op"))
			return
		}

		ctx := context.Background() // request ctx

		// init Firestore client
		if err := h.google.InitFirestore(ctx); err != nil {
			c.Next(err)
			return
		}

		switch *payload.Type {
		case constant.TypeCategory:
			if err := h.storeCategories(ctx, payload); err != nil {
				c.Next(err)
				return
			}
		case constant.TypeFeed:
			if err := h.storeFeed(ctx, payload); err != nil {
				c.Next(err)
				return
			}
		case constant.TypeEntry:
			if err := h.storeEntry(ctx, payload); err != nil {
				c.Next(err)
				return
			}
		}
		c.SendStatus(http.StatusOK)
	}
}
