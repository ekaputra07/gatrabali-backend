package function

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"

	"github.com/apps4bali/gatrabali-backend/common/model"
)

// FirestoreEvent is the payload of a Firestore event.
type FirestoreEvent struct {
	OldValue   FirestoreValue `json:"oldValue"`
	Value      FirestoreValue `json:"value"`
	UpdateMask struct {
		FieldPaths []string `json:"fieldPaths"`
	} `json:"updateMask"`
}

// FirestoreValue holds Firestore fields.
type FirestoreValue struct {
	CreateTime time.Time `json:"createTime"`
	// Fields is the data for this value. The type depends on the format of your
	// database. Log an interface{} value and inspect the result to see a JSON
	// representation of your database fields.
	Fields     model.Entry `json:"fields"`
	Name       string      `json:"name"`
	UpdateTime time.Time   `json:"updateTime"`
}

// NotifyCategorySubscribers triggered when new entry written to Firestore,
// get the list of the subscribers for the category of this entry and send a message to PushNotification topic.
func NotifyCategorySubscribers(ctx context.Context, e FirestoreEvent) error {
	client, err := firebaseApp.FirestoreClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	fmt.Printf("NotifyCategorySubscribers triggered by entry=%v, with categories=%v", e.Value.Fields.ID, e.Value.Fields.Categories)
	categories := e.Value.Fields.Categories
	// if categories set, do nothing
	if len(categories) == 0 {
		return nil
	}

	categoryID := categories[0]

	// Get the category
	doc, err := client.Doc(fmt.Sprintf("/categories/%v", categoryID)).Get(ctx)
	if err != nil {
		fmt.Printf("Category with ID=%v does not exists", categoryID)
		return nil
	}
	category := doc.Data()

	// create message to publish to PushNotification topic.
	pushData := model.PushNotificationPayload{
		Title: fmt.Sprintf("Berita terbaru di %v", category["title"]),
		Body:  e.Value.Fields.Title,
	}

	// create PubSub client
	pubsubClient, err := firebaseApp.PubSubClient(ctx)
	if err != nil {
		return nil
	}
	pushTopic := pubsubClient.Topic("PushNotification")

	// get subscribers
	iter := client.Collection(fmt.Sprintf("/categories/%v/subscribers", categoryID)).Documents(ctx)
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
		pushData.UserID = fmt.Sprintf("%v", subscriber["user_id"]) // set the UserID

		j, _ := json.Marshal(pushData)
		pubsubMsg := &pubsub.Message{Data: j}
		pushTopic.Publish(ctx, pubsubMsg)
	}
	return nil
}
