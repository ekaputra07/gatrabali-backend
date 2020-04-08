package events

import (
	"context"
	"encoding/json"
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
