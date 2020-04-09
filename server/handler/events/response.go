package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"

	"server/common/constant"
	"server/common/service"
	"server/common/types"
	"server/config"
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

			// notify parent comment author
			if (incrementValue > 0) && (parent != nil) && (parent.Data()["user_id"].(string) != r.UserID) {
				if err := r.notifyParentAuthor(ctx, parent.Data()["user_id"].(string)); err != nil {
					log.Printf("[ERROR] %s\n", err)
				}
			}
		}

		// -- transaction end
		return nil
	})
}

func (r *response) notifyParentAuthor(ctx context.Context, parentAuthorID string) error {
	payload := types.PushNotificationPayload{
		UserID: parentAuthorID, // to user
		Title:  fmt.Sprintf("%s membalas komentar anda:", r.User.Name),
		Body:   r.Comment,
		Image:  r.User.Avatar,
		Data: map[string]string{
			"click_action": "FLUTTER_NOTIFICATION_CLICK",
			"data_type":    "response",
			"entry_title":  r.Entry.Title,
			"entry_id":     strconv.FormatInt(r.Entry.ID, 10),
			"category_id":  strconv.FormatInt(r.Entry.CategoryID, 10),
			"feed_id":      strconv.FormatInt(r.Entry.FeedID, 10),
			"published_at": strconv.FormatInt(r.Entry.PublishedAt, 10),
		},
	}
	j, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	serverID, err := r.google.PublishToTopic(ctx, config.PushNotificationTopic, &pubsub.Message{Data: j})
	if err != nil {
		log.Println("notifyParentAuthor(): publish to Push topic failed:", err)
	} else {
		log.Println("notifyParentAuthor(): publish to Push topic success:", serverID)
	}
	return err
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
