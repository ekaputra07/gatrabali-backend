package function

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"

	"github.com/apps4bali/gatrabali-backend/common/constant"
	"github.com/apps4bali/gatrabali-backend/common/model"
)

// EntryEvent is the payload of a Firestore event.
type EntryEvent struct {
	OldValue   EntryValue `json:"oldValue"`
	Value      EntryValue `json:"value"`
	UpdateMask struct {
		FieldPaths []string `json:"fieldPaths"`
	} `json:"updateMask"`
}

// EntryValue is the values in Firestore event.
type EntryValue struct {
	CreateTime time.Time   `json:"createTime"`
	Fields     EntryFields `json:"fields"`
	Name       string      `json:"name"`
	UpdateTime time.Time   `json:"updateTime"`
}

// EntryFields is the entry fields itself, we only need id, title and categories here.
type EntryFields struct {
	ID struct {
		Value string `json:"integerValue"`
	} `json:"id"`
	Title struct {
		Value string `json:"stringValue"`
	} `json:"title"`
	FeedID struct {
		Value string `json:"integerValue"`
	} `json:"feed_id"`
	CategoryID struct {
		Value string `json:"integerValue"`
	} `json:"category_id"`
	PublishedAt struct {
		Value string `json:"integerValue"`
	} `json:"published_at"`
}

// isUserExists check to see if user with given ID is currently exists.
func isUserExists(ctx context.Context, client *firestore.Client, userID string) bool {
	_, err := client.Collection("users").Doc(userID).Get(ctx)
	return err == nil
}

// NotifyCategorySubscribers triggered when new entry written to Firestore,
// get the list of the subscribers for the category of this entry and send a message to PushNotification topic.
func NotifyCategorySubscribers(ctx context.Context, e EntryEvent) error {

	fmt.Printf("NotifyCategorySubscribers triggered by entry=%v, with categories=%v\n",
		e.Value.Fields.ID.Value, e.Value.Fields.CategoryID.Value)

	client, err := firebaseApp.FirestoreClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	entryTitle := e.Value.Fields.Title.Value
	categoryID := e.Value.Fields.CategoryID.Value
	feedID := e.Value.Fields.FeedID.Value

	// overrides categoryID if feedID belongs to BaleBengong
	baleBengongFeeds := []string{"33", "34", "35", "36", "37", "38", "39", "40"}
	for _, ID := range baleBengongFeeds {
		if ID == feedID {
			categoryID = "balebengong"
			break
		}
	}

	// Get the category
	doc, err := client.Collection(constant.Categories).Doc(categoryID).Get(ctx)
	if err != nil {
		fmt.Printf("Category with ID=%v does not exists\n", categoryID)
		return nil
	}
	category := doc.Data()

	// create PubSub client
	pubsubClient, err := firebaseApp.PubSubClient(ctx)
	if err != nil {
		return err
	}
	defer pubsubClient.Close()
	pushTopic := pubsubClient.Topic("PushNotification")

	// create message to publish to PushNotification topic.
	pushData := model.PushNotificationPayload{
		Title: fmt.Sprintf("%v", category["title"]),
		Body:  entryTitle,
		Data: map[string]string{
			"click_action":   "FLUTTER_NOTIFICATION_CLICK",
			"data_type":      "entry",
			"entry_title":    entryTitle,
			"entry_id":       e.Value.Fields.ID.Value,
			"category_id":    e.Value.Fields.CategoryID.Value,
			"category_title": fmt.Sprintf("%v", category["title"]),
			"feed_id":        e.Value.Fields.FeedID.Value,
			"published_at":   e.Value.Fields.PublishedAt.Value,
		},
	}

	// get subscribers
	iter := client.Collection(fmt.Sprintf("categories/%v/subscribers", categoryID)).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}

		subscriber := doc.Data()
		// set recipient
		pushData.UserID = fmt.Sprintf("%v", subscriber["user_id"])

		// check to see if user exists before publishing a message.
		// if user does not exists, delete them from subscriber list.
		if !isUserExists(ctx, client, pushData.UserID) {
			_, err := doc.Ref.Delete(ctx)
			if err != nil {
				fmt.Printf("Failed to delete subscriber %v from category %v\n", pushData.UserID, categoryID)
			}
			continue
		}

		j, _ := json.Marshal(pushData)
		pubsubMsg := &pubsub.Message{Data: j}
		serverID, err := pushTopic.Publish(ctx, pubsubMsg).Get(ctx)
		if err != nil {
			fmt.Printf("Publish to Topic failed: %s\n", err)
		} else {
			fmt.Printf("Publish to Topic success: %s\n", serverID)
		}
	}
	return nil
}
