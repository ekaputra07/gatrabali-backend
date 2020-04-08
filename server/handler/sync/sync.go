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

// Handler sync data from Miniflux to Firestore
func Handler(ctx context.Context, google *service.Google) func(*fiber.Ctx) {
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

		// init Firestore client
		if err := google.InitFirestore(ctx); err != nil {
			c.Next(err)
			return
		}

		ctx := context.Background() // request ctx

		switch *payload.Type {
		case constant.TypeCategory:
			if err := storeCategories(ctx, google.Firestore, payload); err != nil {
				c.Next(err)
				return
			}
		case constant.TypeFeed:
			if err := storeFeed(ctx, google.Firestore, payload); err != nil {
				c.Next(err)
				return
			}
		case constant.TypeEntry:
			if err := storeEntry(ctx, google.Firestore, payload); err != nil {
				c.Next(err)
				return
			}
		}
		c.SendStatus(http.StatusOK)
	}
}
