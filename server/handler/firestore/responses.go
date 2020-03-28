package firestore

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/apps4bali/gatrabali-backend/common/constant"
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

type response struct {
	UserID          string `json:"user_id"`
	Type            string `json:"type"`
	EntryID         int64  `json:"entry_id"`
	EntryCategoryID int64  `json:"entry_category_id"`
	EntryFeedID     int64  `json:"entry_feed_id"`
	ParentID        string `json:"parent_id"`
	Reaction        string `json:"reaction"`
	Comment         string `json:"comment"`
}

func aggregateEntryResponse(ctx context.Context, store *firestore.Client, rawdata []byte) error {
	var data *responseData
	if err := json.Unmarshal(rawdata, &data); err != nil {
		return err
	}

	// on created
	if (data.Before == nil) && (data.After != nil) {
		switch data.After.Type {
		case typeComment:
			return aggregateCommentCreate(ctx, store, data.After)
		case typeReaction:
			return aggregateReactionCreate(ctx, store, data.After)
		}
	}

	// on updated
	// only process REACTION update, comment update does not affect comments count.
	if (data.Before != nil) && (data.After != nil) {
		if data.After.Type == typeReaction {
			return aggregateReactionUpdate(ctx, store, data.Before, data.After)
		}
	}

	// on deleted
	if (data.Before != nil) && (data.After == nil) {
		switch data.Before.Type {
		case typeComment:
			return aggregateCommentDelete(ctx, store, data.Before)
		case typeReaction:
			return aggregateReactionDelete(ctx, store, data.Before)
		}
	}

	return nil
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

// -- comment aggregation
func aggregateCommentCreate(ctx context.Context, client *firestore.Client, resp *response) error {
	entryID := strconv.FormatInt(resp.EntryID, 10)
	categoryID := resp.EntryCategoryID

	update := []firestore.Update{{
		Path:  "comment_count",
		Value: firestore.Increment(1),
	}}
	_, err := client.Collection(entryCollectionByCategory(categoryID)).Doc(entryID).Update(ctx, update)
	return err
}

func aggregateCommentDelete(ctx context.Context, client *firestore.Client, resp *response) error {
	entryID := strconv.FormatInt(resp.EntryID, 10)
	categoryID := resp.EntryCategoryID

	update := []firestore.Update{{
		Path:  "comment_count",
		Value: firestore.Increment(-1),
	}}
	_, err := client.Collection(entryCollectionByCategory(categoryID)).Doc(entryID).Update(ctx, update)
	return err
}

// -- reaction aggregation
func aggregateReactionCreate(ctx context.Context, client *firestore.Client, resp *response) error {
	entryID := strconv.FormatInt(resp.EntryID, 10)
	categoryID := resp.EntryCategoryID
	reaction := resp.Reaction

	update := []firestore.Update{{
		Path:  fmt.Sprintf("reaction_%s_count", strings.ToLower(reaction)),
		Value: firestore.Increment(1),
	}}
	_, err := client.Collection(entryCollectionByCategory(categoryID)).Doc(entryID).Update(ctx, update)
	return err
}

func aggregateReactionUpdate(ctx context.Context, client *firestore.Client, before, after *response) error {
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
		{Path: fmt.Sprintf("reaction_%s_count", strings.ToLower(newReaction)),
			Value: firestore.Increment(1),
		},
	}
	_, err := client.Collection(entryCollectionByCategory(categoryID)).Doc(entryID).Update(ctx, update)
	return err
}

func aggregateReactionDelete(ctx context.Context, client *firestore.Client, resp *response) error {
	entryID := strconv.FormatInt(resp.EntryID, 10)
	categoryID := resp.EntryCategoryID
	reaction := resp.Reaction

	update := []firestore.Update{{
		Path:  fmt.Sprintf("reaction_%s_count", strings.ToLower(reaction)),
		Value: firestore.Increment(-1),
	}}
	_, err := client.Collection(entryCollectionByCategory(categoryID)).Doc(entryID).Update(ctx, update)
	return err
}
