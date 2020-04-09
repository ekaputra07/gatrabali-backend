package events

import (
	"context"
	"encoding/json"
	"fmt"
	"server/common/constant"
	"strconv"
	"strings"

	"cloud.google.com/go/firestore"
)

const (
	typeComment  = "COMMENT"
	typeReaction = "REACTION"
)

// this is based PubSub data format sent by Firesub
// https://github.com/ekaputra07/firesub
type responseData struct {
	ID        string    `json:"id"`
	Timestamp string    `json:"timestamp"`
	Before    *response `json:"before"`
	After     *response `json:"after"`
}

func (h *Handler) aggregateResponses(ctx context.Context, pubsubData []byte) error {
	var data *responseData
	if err := json.Unmarshal(pubsubData, &data); err != nil {
		return err
	}

	// on created
	if (data.Before == nil) && (data.After != nil) {
		after := data.After.setGoogle(h.google)

		switch after.Type {
		case typeComment:
			return after.aggregateComment(ctx, 1)
		case typeReaction:
			return after.aggregateReactionCreateDelete(ctx, 1)
		}
	}

	// on updated
	// only process REACTION update, comment update does not affect comments count.
	if (data.Before != nil) && (data.After != nil) {
		after := data.After.setGoogle(h.google)

		if after.Type == typeReaction {
			return aggregateReactionUpdate(ctx, data.Before, after)
		}
	}

	// on deleted
	if (data.Before != nil) && (data.After == nil) {
		before := data.Before.setGoogle(h.google)

		switch before.Type {
		case typeComment:
			// aggregator
			err := before.aggregateComment(ctx, -1)
			if err != nil {
				return err
			}
			// delete replies if any
			err = before.deleteReplies(ctx, data.ID)
			if err != nil {
				return err
			}
		case typeReaction:
			return before.aggregateReactionCreateDelete(ctx, -1)
		}
	}

	return nil
}

func aggregateReactionUpdate(ctx context.Context, before, after *response) error {
	entryID := strconv.FormatInt(after.EntryID, 10)
	categoryID := after.EntryCategoryID
	newReaction := after.Reaction
	oldReaction := before.Reaction

	if newReaction == oldReaction {
		return nil
	}

	update := []firestore.Update{
		{
			Path:  fmt.Sprintf("reaction_%s_count", strings.ToLower(oldReaction)),
			Value: firestore.Increment(-1),
		},
		{
			Path:  fmt.Sprintf("reaction_%s_count", strings.ToLower(newReaction)),
			Value: firestore.Increment(1),
		},
	}
	_, err := after.google.Firestore.
		Collection(entryCollectionByCategory(categoryID)).
		Doc(entryID).
		Update(ctx, update)

	return err
}

// entryCollectionByCategory method is specific to BaliFeed app only
func entryCollectionByCategory(categoryID int64) string {
	c := constant.Entries
	if categoryID == 11 {
		c = constant.Kriminal
	} else if categoryID == 12 {
		c = constant.BaliUnited
	} else if categoryID > 12 {
		c = constant.BaleBengong
	}
	return c
}
