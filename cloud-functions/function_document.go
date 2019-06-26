package function

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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

	// categories:map[arrayValue:map[values:[map[integerValue:6]]]]
	Categories struct {
		Value struct {
			Values []struct {
				Value string `json:"integerValue"`
			} `json:"values"`
		} `json:"arrayValue"`
	} `json:"categories"`
}

// NotifyCategorySubscribers triggered when new entry written to Firestore,
// get the list of the subscribers for the category of this entry and send a message to PushNotification topic.
func NotifyCategorySubscribers(ctx context.Context, e EntryEvent) error {

	fmt.Printf("NotifyCategorySubscribers triggered by entry=%v, with categories=%v\n",
		e.Value.Fields.ID.Value, e.Value.Fields.Categories.Value)

	client, err := firebaseApp.FirestoreClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	entryTitle := e.Value.Fields.Title.Value
	categories := e.Value.Fields.Categories.Value.Values
	// if categories set, do nothing
	if len(categories) == 0 {
		return nil
	}

	categoryID := categories[0].Value

	// Get the category
	doc, err := client.Collection(constant.Categories).Doc(fmt.Sprintf("%v", categoryID)).Get(ctx)
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
		Title: fmt.Sprintf("Gatra %v", category["title"]),
		Body:  entryTitle,
		Data:  map[string]string{"uri": fmt.Sprintf("gatrabali://entries/%v", e.Value.Fields.ID.Value)},
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
