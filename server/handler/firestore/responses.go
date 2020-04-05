package firestore

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"cloud.google.com/go/firestore"

	"server/common/constant"
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

type user struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type response struct {
	UserID          string `json:"user_id"`
	Type            string `json:"type"`
	EntryID         int64  `json:"entry_id"`
	EntryCategoryID int64  `json:"entry_category_id"`
	EntryFeedID     int64  `json:"entry_feed_id"`
	ParentID        string `json:"parent_id"` // level 0 or level n
	ThreadID        string `json:"thread_id"` // level 0
	Reaction        string `json:"reaction"`
	Comment         string `json:"comment"`
	Entry           entry  `json:"entry"`
	User            user   `json:"user"`
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
			return aggregateComment(ctx, store, data.After, 1)
		case typeReaction:
			return aggregateReactionCreateDelete(ctx, store, data.After, 1)
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
			// aggregator
			err := aggregateComment(ctx, store, data.Before, -1)
			if err != nil {
				return err
			}
			// delete replies if any
			err = deleteCommentReplies(ctx, store, data.ID, data.Before)
			if err != nil {
				return err
			}
		case typeReaction:
			return aggregateReactionCreateDelete(ctx, store, data.Before, -1)
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

// deleteCommentReplies deletes all replies for this comment
func deleteCommentReplies(ctx context.Context, client *firestore.Client, ID string, resp *response) error {
	// top level comment, delete all replies
	if resp.ThreadID == "" {
		iter := client.Collection(constant.EntryResponses).Where("thread_id", "==", ID).Documents(ctx)
		snaps, err := iter.GetAll()
		if err != nil {
			return err
		}
		if len(snaps) > 0 {
			batch := client.Batch()
			for _, snap := range snaps {
				batch.Delete(snap.Ref)
			}
			_, err := batch.Commit(ctx)
			return err
		}
	}

	// a reply, delete all replies (childs) to this reply
	if resp.ThreadID != "" {
		iter := client.Collection(constant.EntryResponses).Where("parent_id", "==", ID).Documents(ctx)
		snaps, err := iter.GetAll()
		if err != nil {
			return err
		}
		if len(snaps) > 0 {
			batch := client.Batch()
			for _, snap := range snaps {
				batch.Delete(snap.Ref)
			}
			_, err := batch.Commit(ctx)
			return err
		}
	}

	return nil
}

// -- comment aggregation
func aggregateComment(ctx context.Context, client *firestore.Client, resp *response, incrementValue int) error {
	entryID := strconv.FormatInt(resp.EntryID, 10)
	categoryID := resp.EntryCategoryID

	// run inside a transaction
	return client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// -- transaction start

		// update entry comment count
		entry := client.Collection(entryCollectionByCategory(categoryID)).Doc(entryID)
		update := []firestore.Update{{
			Path:  "comment_count",
			Value: firestore.Increment(incrementValue),
		}}
		if err := tx.Update(entry, update); err != nil {
			return err
		}

		// if has parent_id and thread_id
		if resp.ParentID != "" && resp.ThreadID != "" {

			var parent *firestore.DocumentSnapshot
			var thread *firestore.DocumentSnapshot
			var err error

			// get direct parent
			parent, _ = client.Collection(constant.EntryResponses).Doc(resp.ParentID).Get(ctx)

			// a reply to a reply, get level 0 parent (thread)
			if resp.ParentID != resp.ThreadID {
				thread, _ = client.Collection(constant.EntryResponses).Doc(resp.ThreadID).Get(ctx)
			}

			// if thread found, increment reply_count
			// otherwise increment reply_count of parent
			update := []firestore.Update{{
				Path:  "reply_count",
				Value: firestore.Increment(incrementValue),
			}}
			// update reply count on parent
			if parent != nil {
				err = tx.Update(parent.Ref, update)
				if err != nil {
					return err
				}
			}
			// update reply count on thread
			if thread != nil {
				err = tx.Update(thread.Ref, update)
				if err != nil {
					return err
				}
			}

			// TODO: Notify parent author
			// if failed should not failed the transaction.
			if (incrementValue > 0) && (parent != nil) && (resp.UserID != parent.Data()["user_id"].(string)) {
				fmt.Println("TODO: send push!")
			}
		}

		// -- transaction end
		return nil
	})
}

// -- reaction aggregation
func aggregateReactionCreateDelete(ctx context.Context, client *firestore.Client, resp *response, incrementValue int) error {
	entryID := strconv.FormatInt(resp.EntryID, 10)
	categoryID := resp.EntryCategoryID
	reaction := resp.Reaction

	update := []firestore.Update{{
		Path:  fmt.Sprintf("reaction_%s_count", strings.ToLower(reaction)),
		Value: firestore.Increment(incrementValue),
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
		{
			Path:  fmt.Sprintf("reaction_%s_count", strings.ToLower(newReaction)),
			Value: firestore.Increment(1),
		},
	}
	_, err := client.Collection(entryCollectionByCategory(categoryID)).Doc(entryID).Update(ctx, update)
	return err
}
