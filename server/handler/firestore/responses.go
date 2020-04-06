package firestore

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"

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

func (h *Handler) aggregateResponses(ctx context.Context, pubsubData []byte) error {
	var data *responseData
	if err := json.Unmarshal(pubsubData, &data); err != nil {
		return err
	}

	// on created
	if (data.Before == nil) && (data.After != nil) {
		after := data.After.setFirestore(h.firestore)

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
		after := data.After.setFirestore(h.firestore)

		if after.Type == typeReaction {
			return aggregateReactionUpdate(ctx, data.Before, after)
		}
	}

	// on deleted
	if (data.Before != nil) && (data.After == nil) {
		before := data.Before.setFirestore(h.firestore)

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

type user struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type response struct {
	firestore *firestore.Client
	pubsub    *pubsub.Client

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

func (r *response) setFirestore(client *firestore.Client) *response {
	r.firestore = client
	return r
}

func (r *response) setPubsub(client *pubsub.Client) *response {
	r.pubsub = client
	return r
}

// deleteReplies deletes all replies for this comment
func (r *response) deleteReplies(ctx context.Context, ID string) error {
	// top level comment, delete all replies
	if r.ThreadID == "" {
		iter := r.firestore.Collection(constant.EntryResponses).Where("thread_id", "==", ID).Documents(ctx)
		snaps, err := iter.GetAll()
		if err != nil {
			return err
		}
		if len(snaps) > 0 {
			batch := r.firestore.Batch()
			for _, snap := range snaps {
				batch.Delete(snap.Ref)
			}
			_, err := batch.Commit(ctx)
			return err
		}
	}

	// a reply, delete all replies (childs) to this reply
	if r.ThreadID != "" {
		iter := r.firestore.Collection(constant.EntryResponses).Where("parent_id", "==", ID).Documents(ctx)
		snaps, err := iter.GetAll()
		if err != nil {
			return err
		}
		if len(snaps) > 0 {
			batch := r.firestore.Batch()
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
func (r *response) aggregateComment(ctx context.Context, incrementValue int) error {
	entryID := strconv.FormatInt(r.EntryID, 10)
	categoryID := r.EntryCategoryID

	// run inside a transaction
	return r.firestore.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// -- transaction start

		// update entry comment count
		entry := r.firestore.Collection(entryCollectionByCategory(categoryID)).Doc(entryID)
		update := []firestore.Update{{
			Path:  "comment_count",
			Value: firestore.Increment(incrementValue),
		}}
		if err := tx.Update(entry, update); err != nil {
			return err
		}

		// if has parent_id and thread_id
		if r.ParentID != "" && r.ThreadID != "" {

			var parent *firestore.DocumentSnapshot
			var thread *firestore.DocumentSnapshot
			var err error

			// get direct parent
			parent, err = r.firestore.Collection(constant.EntryResponses).Doc(r.ParentID).Get(ctx)
			if err != nil {
				parent = nil
			}

			// a reply to a reply, get level 0 parent (thread)
			if r.ParentID != r.ThreadID {
				thread, _ = r.firestore.Collection(constant.EntryResponses).Doc(r.ThreadID).Get(ctx)
				if err != nil {
					thread = nil
				}
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
			if (incrementValue > 0) && (parent != nil) && (r.UserID != parent.Data()["user_id"].(string)) {
				fmt.Println("TODO: send push!")
			}
		}

		// -- transaction end
		return nil
	})
}

// -- reaction aggregation
func (r *response) aggregateReactionCreateDelete(ctx context.Context, incrementValue int) error {
	entryID := strconv.FormatInt(r.EntryID, 10)
	categoryID := r.EntryCategoryID
	reaction := r.Reaction

	update := []firestore.Update{{
		Path:  fmt.Sprintf("reaction_%s_count", strings.ToLower(reaction)),
		Value: firestore.Increment(incrementValue),
	}}
	_, err := r.firestore.
		Collection(entryCollectionByCategory(categoryID)).
		Doc(entryID).
		Update(ctx, update)

	return err
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
	_, err := after.firestore.
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
