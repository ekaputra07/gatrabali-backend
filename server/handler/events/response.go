package events

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"cloud.google.com/go/firestore"

	"server/common/constant"
	"server/common/service"
)

type user struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type response struct {
	google *service.Google

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

func (r *response) setGoogle(g *service.Google) *response {
	r.google = g
	return r
}

// deleteReplies deletes all replies for this comment
func (r *response) deleteReplies(ctx context.Context, ID string) error {
	// top level comment, delete all replies
	if r.ThreadID == "" {
		iter := r.google.Firestore.Collection(constant.EntryResponses).Where("thread_id", "==", ID).Documents(ctx)
		snaps, err := iter.GetAll()
		if err != nil {
			return err
		}
		if len(snaps) > 0 {
			batch := r.google.Firestore.Batch()
			for _, snap := range snaps {
				batch.Delete(snap.Ref)
			}
			_, err := batch.Commit(ctx)
			return err
		}
	}

	// a reply, delete all replies (childs) to this reply
	if r.ThreadID != "" {
		iter := r.google.Firestore.Collection(constant.EntryResponses).Where("parent_id", "==", ID).Documents(ctx)
		snaps, err := iter.GetAll()
		if err != nil {
			return err
		}
		if len(snaps) > 0 {
			batch := r.google.Firestore.Batch()
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
	return r.google.Firestore.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// -- transaction start

		// update entry comment count
		entry := r.google.Firestore.Collection(entryCollectionByCategory(categoryID)).Doc(entryID)
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
			parent, err = r.google.Firestore.Collection(constant.EntryResponses).Doc(r.ParentID).Get(ctx)
			if err != nil {
				parent = nil
			}

			// a reply to a reply, get level 0 parent (thread)
			if r.ParentID != r.ThreadID {
				thread, _ = r.google.Firestore.Collection(constant.EntryResponses).Doc(r.ThreadID).Get(ctx)
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
	_, err := r.google.Firestore.
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
