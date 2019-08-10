package function

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/apps4bali/gatrabali-backend/common/constant"
)

// responseEvent is the payload of a Firestore event.
type responseEvent struct {
	OldValue   responseValue `json:"oldValue"`
	Value      responseValue `json:"value"`
	UpdateMask struct {
		FieldPaths []string `json:"fieldPaths"`
	} `json:"updateMask"`
}

// responseValue is the values in Firestore event.
type responseValue struct {
	CreateTime time.Time      `json:"createTime"`
	Fields     responseFields `json:"fields"`
	Name       string         `json:"name"`
	UpdateTime time.Time      `json:"updateTime"`
}

// responseFields is the response fields itself, we only need id, title and categories here.
type responseFields struct {
	UserID struct {
		Value string `json:"stringValue"`
	} `json:"user_id"`
	Type struct {
		Value string `json:"stringValue"`
	} `json:"type"`
	EntryID struct {
		Value string `json:"integerValue"`
	} `json:"entry_id"`
	EntryCategoryID struct {
		Value string `json:"integerValue"`
	} `json:"entry_category_id"`
	EntryFeedID struct {
		Value string `json:"integerValue"`
	} `json:"entry_feed_id"`
	ParentID struct {
		Value string `json:"stringValue"`
	} `json:"parent_id"`
	Reaction struct {
		Value string `json:"stringValue"`
	} `json:"reaction"`
	Comment struct {
		Value string `json:"stringValue"`
	} `json:"comment"`
}

// AggregateEntryResponses triggered when new entry response created, updated or deleted on Firestore,
// It aggregate responses count based on its type and update the entry.
func AggregateEntryResponses(ctx context.Context, e responseEvent) error {
	fmt.Printf("Triggered by a response in entry=%v\n", e.Value.Fields.EntryID.Value)

	client, err := firebaseApp.FirestoreClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	// after created
	if e.Value.Name != "" && e.OldValue.Name == "" {
		respType := e.Value.Fields.Type.Value
		switch respType {
		case "COMMENT":
			return aggregateCommentCreate(ctx, client, e.Value.Fields)
		case "REACTION":
			return aggregateReactionCreate(ctx, client, e.Value.Fields)
		}
	}

	// after updated
	if e.Value.Name != "" && e.OldValue.Name != "" {
		respType := e.Value.Fields.Type.Value
		// only process REACTION update, comment update does not affect comments count.
		if respType == "REACTION" {
			return aggregateReactionUpdate(ctx, client, e.OldValue.Fields, e.Value.Fields)
		}
	}

	// after deleted
	if e.Value.Name == "" && e.OldValue.Name != "" {
		respType := e.OldValue.Fields.Type.Value
		switch respType {
		case "COMMENT":
			return aggregateCommentDelete(ctx, client, e.OldValue.Fields)
		case "REACTION":
			return aggregateReactionDelete(ctx, client, e.OldValue.Fields)
		}
	}
	return nil
}

func entryCollectionByCategory(categoryID int) string {
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
func aggregateCommentCreate(ctx context.Context, client *firestore.Client, fields responseFields) error {
	entryID := fields.EntryID.Value
	categoryID, _ := strconv.Atoi(fields.EntryCategoryID.Value)

	update := []firestore.Update{{Path: "comment_count", Value: firestore.Increment(1)}}
	_, err := client.Collection(entryCollectionByCategory(categoryID)).Doc(entryID).Update(ctx, update)
	return err
}

func aggregateCommentDelete(ctx context.Context, client *firestore.Client, oldFields responseFields) error {
	entryID := oldFields.EntryID.Value
	categoryID, _ := strconv.Atoi(oldFields.EntryCategoryID.Value)

	update := []firestore.Update{{Path: "comment_count", Value: firestore.Increment(-1)}}
	_, err := client.Collection(entryCollectionByCategory(categoryID)).Doc(entryID).Update(ctx, update)
	return err
}

// -- reaction aggregation
func aggregateReactionCreate(ctx context.Context, client *firestore.Client, fields responseFields) error {
	entryID := fields.EntryID.Value
	categoryID, _ := strconv.Atoi(fields.EntryCategoryID.Value)
	reaction := fields.Reaction.Value

	update := []firestore.Update{{Path: fmt.Sprintf("reaction_%s_count", strings.ToLower(reaction)), Value: firestore.Increment(1)}}
	_, err := client.Collection(entryCollectionByCategory(categoryID)).Doc(entryID).Update(ctx, update)
	return err
}

func aggregateReactionUpdate(ctx context.Context, client *firestore.Client, oldFields, newFields responseFields) error {
	entryID := newFields.EntryID.Value
	categoryID, _ := strconv.Atoi(newFields.EntryCategoryID.Value)
	newReaction := newFields.Reaction.Value
	oldReaction := oldFields.Reaction.Value

	if newReaction == oldReaction {
		return nil
	}

	update := []firestore.Update{
		{Path: fmt.Sprintf("reaction_%s_count", strings.ToLower(oldReaction)), Value: firestore.Increment(-1)},
		{Path: fmt.Sprintf("reaction_%s_count", strings.ToLower(newReaction)), Value: firestore.Increment(1)},
	}
	_, err := client.Collection(entryCollectionByCategory(categoryID)).Doc(entryID).Update(ctx, update)
	return err
}

func aggregateReactionDelete(ctx context.Context, client *firestore.Client, oldFields responseFields) error {
	entryID := oldFields.EntryID.Value
	categoryID, _ := strconv.Atoi(oldFields.EntryCategoryID.Value)
	oldReaction := oldFields.Reaction.Value

	update := []firestore.Update{{Path: fmt.Sprintf("reaction_%s_count", strings.ToLower(oldReaction)), Value: firestore.Increment(-1)}}
	_, err := client.Collection(entryCollectionByCategory(categoryID)).Doc(entryID).Update(ctx, update)
	return err
}
