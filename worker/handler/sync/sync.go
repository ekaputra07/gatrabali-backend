package sync

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/apps4bali/gatrabali-backend/common/constant"
	"github.com/apps4bali/gatrabali-backend/common/model"
	"github.com/gofiber/fiber"

	"worker/handler/sync/service"
)

// Handler sync data from Miniflux to Firestore
func Handler(client *firestore.Client) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		ctx := context.Background()
		data := c.Locals("PubSubData").([]byte)

		log.Printf("Sync triggered with payload: %v\n", string(data))

		var payload *model.SyncPayload
		if err := json.Unmarshal(data, &payload); err != nil {
			c.Next(err)
			return
		}
		if payload.ID == nil || payload.Type == nil || payload.Op == nil {
			c.Next(errors.New("Invalid message payload: missing id, type or op"))
			return
		}

		switch *payload.Type {
		case constant.TypeCategory:
			if err := service.StartCategorySync(ctx, client, payload); err != nil {
				c.Next(err)
			}
		case constant.TypeFeed:
			if err := service.StartFeedSync(ctx, client, payload); err != nil {
				c.Next(err)
			}
		case constant.TypeEntry:
			if err := service.StartEntrySync(ctx, client, payload); err != nil {
				c.Next(err)
			}
		}
		c.SendStatus(http.StatusOK)
	}
}
