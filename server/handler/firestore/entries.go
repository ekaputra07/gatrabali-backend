package firestore

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"github.com/apps4bali/gatrabali-backend/common/constant"
	"github.com/apps4bali/gatrabali-backend/common/types"
	"google.golang.org/api/iterator"

	"server/config"
)

// this is based PubSub data formate sent by Firesub
// https://github.com/ekaputra07/firesub
type entryData struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Entry     *entry `json:"data"`
}

type entry struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	FeedID      int64  `json:"feed_id"`
	CategoryID  int64  `json:"category_id"`
	PublishedAt int64  `json:"published_at"`
}

// isUserExists check to see if user with given ID is currently exists.
func isUserExists(ctx context.Context, client *firestore.Client, userID string) bool {
	_, err := client.Collection("users").Doc(userID).Get(ctx)
	return err == nil
}

// notifySubscriber triggered when new entry written to Firestore,
// get the list of the subscribers for the category of this entry and send a message to PushNotification topic.
// TODO: this method is too long!
func notifySubscriber(
	ctx context.Context,
	firestoreClient *firestore.Client,
	pubsubClient *pubsub.Client,
	rawdata []byte) error {

	var data *entryData
	if err := json.Unmarshal(rawdata, &data); err != nil {
		return err
	}

	entryTitle := data.Entry.Title
	entryID := strconv.FormatInt(data.Entry.ID, 10)
	categoryID := strconv.FormatInt(data.Entry.CategoryID, 10)
	feedID := strconv.FormatInt(data.Entry.FeedID, 10)
	publishedAt := strconv.FormatInt(data.Entry.PublishedAt, 10)

	// overrides categoryID if feedID belongs to BaleBengong
	baleBengongFeeds := []string{"33", "34", "35", "36", "37", "38", "39", "40"}
	for _, ID := range baleBengongFeeds {
		if ID == feedID {
			categoryID = "balebengong"
			break
		}
	}

	// Get the category
	doc, err := firestoreClient.Collection(constant.Categories).Doc(categoryID).Get(ctx)
	if err != nil {
		return fmt.Errorf("Category with ID=%v does not exists", categoryID)
	}

	category := doc.Data()
	pushTopic := pubsubClient.Topic(config.PushNotificationTopic)

	// create message to publish to PushNotification topic.
	pushData := types.PushNotificationPayload{
		Title: category["title"].(string),
		Body:  entryTitle,
		Data: map[string]string{
			"click_action":   "FLUTTER_NOTIFICATION_CLICK",
			"data_type":      "entry",
			"entry_title":    entryTitle,
			"entry_id":       entryID,
			"category_id":    categoryID,
			"category_title": category["title"].(string),
			"feed_id":        feedID,
			"published_at":   publishedAt,
		},
	}

	// get subscribers
	iter := firestoreClient.Collection(fmt.Sprintf("categories/%v/subscribers", categoryID)).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
			continue
		}

		subscriber := doc.Data()
		pushData.UserID = subscriber["user_id"].(string) // set recipient

		// check to see if user exists before publishing a message.
		// if user does not exists, delete them from subscriber list.
		if !isUserExists(ctx, firestoreClient, pushData.UserID) {
			if _, err := doc.Ref.Delete(ctx); err != nil {
				log.Printf("Failed to delete subscriber %v from category %v\n", pushData.UserID, categoryID)
			}
			continue
		}

		j, err := json.Marshal(pushData)
		if err != nil {
			log.Println("Failed Marshalling push data:", err)
			continue
		}

		pubsubMsg := &pubsub.Message{Data: j}
		serverID, err := pushTopic.Publish(ctx, pubsubMsg).Get(ctx)
		if err != nil {
			log.Println("Publish to Push topic failed:", err)
		} else {
			log.Println("Publish to Push topic success:", serverID)
		}
	}

	return nil
}
