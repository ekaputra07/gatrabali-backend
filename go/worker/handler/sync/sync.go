package sync

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/apps4bali/gatrabali-backend/go/common"
	"github.com/fiberweb/pubsub"
	"github.com/gofiber/fiber"

	"worker/handler/sync/service"
)

// Handler sync data from Miniflux to Firestore
func Handler(client *firestore.Client) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		ctx := context.Background()
		msg := c.Locals(pubsub.LocalsKey).(*pubsub.Message)

		log.Printf("Sync triggered with payload: %v\n", msg)

		var payload *common.SyncPayload
		if err := json.Unmarshal(msg.Message.Data, &payload); err != nil {
			c.Next(err)
			return
		}
		if payload.ID == nil || payload.Type == nil || payload.Op == nil {
			c.Next(errors.New("Invalid message payload: missing id, type or op"))
			return
		}

		switch *payload.Type {
		case common.TypeCategory:
			if err := service.StartCategorySync(ctx, client, payload); err != nil {
				c.Next(err)
				return
			}
		case common.TypeFeed:
			if err := service.StartFeedSync(ctx, client, payload); err != nil {
				c.Next(err)
				return
			}
		case common.TypeEntry:
			if err := service.StartEntrySync(ctx, client, payload); err != nil {
				c.Next(err)
				return
			}
		}
		c.SendStatus(http.StatusOK)
	}
}
