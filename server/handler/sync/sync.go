package sync

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/apps4bali/gatrabali-backend/common/constant"
	"github.com/apps4bali/gatrabali-backend/common/types"
	"github.com/fiberweb/pubsub"
	"github.com/gofiber/fiber"

	"server/firebase"
	"server/handler/sync/service"
)

// Handler sync data from Miniflux to Firestore
func Handler(ctx context.Context, fb *firebase.Firebase) func(*fiber.Ctx) {
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

		// load Firestore client
		firestore, err := fb.FirestoreClient(ctx)
		if err != nil {
			c.Next(err)
			return
		}

		ctx := context.Background() // request ctx

		switch *payload.Type {
		case constant.TypeCategory:
			if err := service.StartCategorySync(ctx, firestore, payload); err != nil {
				c.Next(err)
				return
			}
		case constant.TypeFeed:
			if err := service.StartFeedSync(ctx, firestore, payload); err != nil {
				c.Next(err)
				return
			}
		case constant.TypeEntry:
			if err := service.StartEntrySync(ctx, firestore, payload); err != nil {
				c.Next(err)
				return
			}
		}
		c.SendStatus(http.StatusOK)
	}
}
