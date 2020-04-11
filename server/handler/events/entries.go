package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"

	"server/common/constant"
	"server/common/types"
	"server/config"
)

// this is based PubSub data formate sent by Firesub
// https://github.com/ekaputra07/firesub
type entryData struct {
	ID        string       `json:"id"`
	Timestamp string       `json:"timestamp"`
	Entry     *types.Entry `json:"data"`
}

// notifySubscriber triggered when new entry written to Firestore,
// get the list of the subscribers for the category of this entry and send a message to PushNotification topic.
// TODO: this method is too long!
func (h *Handler) notifySubscribers(ctx context.Context, pubsubData []byte) error {

	var data *entryData
	if err := json.Unmarshal(pubsubData, &data); err != nil {
		return err
	}

	entryTitle := data.Entry.Title
	entryID := strconv.FormatInt(data.Entry.ID, 10)
	categoryID := strconv.FormatInt(data.Entry.CategoryID, 10)
	feedID := strconv.FormatInt(data.Entry.FeedID, 10)
	publishedAt := strconv.FormatInt(data.Entry.PublishedAt, 10)
	// get entry image
	var entryImage string
	if data.Entry.Enclosures != nil {
		for _, i := range *data.Entry.Enclosures {
			if i.URL != "" {
				entryImage = i.URL
				break
			}
		}
	}

	// overrides subscriberCategory if feedID belongs to BaleBengong
	// for bale bengong subscriber they subscribes on a separate category called "balebengong"
	subscriberCategory := categoryID
	baleBengongFeeds := []string{"33", "34", "35", "36", "37", "38", "39", "40"}
	for _, ID := range baleBengongFeeds {
		if ID == feedID {
			subscriberCategory = "balebengong"
			break
		}
	}

	// Get the category
	doc, err := h.google.Firestore.Collection(constant.Categories).Doc(subscriberCategory).Get(ctx)
	if err != nil {
		return fmt.Errorf("Category with ID=%v does not exists", subscriberCategory)
	}

	category := doc.Data()

	// create message to publish to PushNotification topic.
	pushData := types.PushNotificationPayload{
		Title: category["title"].(string),
		Body:  entryTitle,
		Image: entryImage,
		Data: map[string]string{
			"click_action":   "FLUTTER_NOTIFICATION_CLICK",
			"data_type":      "entry",
			"entry_title":    entryTitle,
			"entry_id":       entryID,
			"category_id":    categoryID, // original category ID
			"category_title": category["title"].(string),
			"feed_id":        feedID,
			"published_at":   publishedAt,
		},
	}

	// get subscribers
	iter := h.google.Firestore.
		Collection(fmt.Sprintf("categories/%v/subscribers", subscriberCategory)).
		Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			continue
		}

		subscriber := doc.Data()
		pushData.UserID = subscriber["user_id"].(string) // set recipient

		// check to see if user exists before publishing a message.
		// if user does not exists, delete them from subscriber list.
		if !h.isUserExists(ctx, pushData.UserID) {
			if _, err := doc.Ref.Delete(ctx); err != nil {
				log.Printf("Failed to delete subscriber %v from category %v\n", pushData.UserID, subscriberCategory)
			}
			continue
		}

		j, err := json.Marshal(pushData)
		if err != nil {
			log.Println("Failed Marshalling push data:", err)
			continue
		}

		pubsubMsg := &pubsub.Message{Data: j}
		if _, err = h.google.PublishToTopic(ctx, config.PushNotificationTopic, pubsubMsg); err != nil {
			log.Println("notifySubscribers(): publish to Push topic failed:", err)
		}
	}

	return nil
}

// h.isUserExists check to see if user with given ID is currently exists.
func (h *Handler) isUserExists(ctx context.Context, userID string) bool {
	_, err := h.google.Firestore.Collection("users").Doc(userID).Get(ctx)
	return err == nil
}
