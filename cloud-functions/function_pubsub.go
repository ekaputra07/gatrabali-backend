package function

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/messaging"

	"github.com/apps4bali/gatrabali-backend/common/constant"
	"github.com/apps4bali/gatrabali-backend/common/model"

	"function/sync"
)

// PubSubMessage is the payload of Pub/Sub message
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// SyncData calls Miniflux API and store its response to Cloud Firestore
func SyncData(ctx context.Context, m PubSubMessage) error {
	log.Printf("SyncData triggered with payload: %v\n", string(m.Data))

	var payload model.SyncPayload
	if err := json.Unmarshal(m.Data, &payload); err != nil {
		return err
	}
	if payload.ID == nil || payload.Type == nil || payload.Op == nil {
		return errors.New("Invalid message payload: missing id, type or op")
	}

	firestore, err := Firestore()
	if err != nil {
		return err
	}
	defer firestore.Close()

	switch *payload.Type {
	case constant.TypeCategory:
		return sync.StartCategorySync(ctx, firestore, &payload)
	case constant.TypeFeed:
		return sync.StartFeedSync(ctx, firestore, &payload)
	case constant.TypeEntry:
		return sync.StartEntrySync(ctx, firestore, &payload)
	}
	return nil
}

// SendPushNotification send push notification using FCM.
func SendPushNotification(ctx context.Context, m PubSubMessage) error {
	log.Printf("SendPushNotification triggered with payload: %v\n", string(m.Data))

	// validate payload
	var payload model.PushNotificationPayload
	if err := json.Unmarshal(m.Data, &payload); err != nil {
		return err
	}
	if payload.Title == "" || payload.Body == "" || payload.UserID == "" {
		return errors.New("Invalid message payload: missing user_id, title or body")
	}

	// get user's FCM tokens
	fclient, err := Firestore()
	if err != nil {
		return err
	}
	defer fclient.Close()

	doc, err := fclient.Collection("users").Doc(payload.UserID).Get(ctx)
	if err != nil {
		return err
	}
	user := doc.Data()
	tokens, ok := user["fcm_tokens"]
	if !ok {
		log.Printf("User %v doesn't have FCM tokens", payload.UserID)
		return nil
	}
	tokensMap := tokens.(map[string]interface{}) // convert to map
	if len(tokensMap) == 0 {
		log.Printf("User %v doesn't have FCM tokens", payload.UserID)
		return nil
	}

	// build notification message
	client, err := MessagingClient(ctx)
	if err != nil {
		return err
	}

	notification := &messaging.Notification{
		Title: payload.Title,
		Body:  payload.Body,
	}

	// loop through tokens and send the notification
	for token := range tokensMap {
		message := &messaging.Message{
			Notification: notification,
			Token:        token,
		}
		resp, err := client.Send(ctx, message)
		if err != nil {
			// if error, delete token
			log.Printf("Notification not sent: %v\n", err)
			delete(tokensMap, token)
		} else {
			log.Printf("Notification sent: %v\n", resp)
		}
	}

	// store back the remaining tokens to user document
	_, err = fclient.Collection("users").Doc(payload.UserID).Update(ctx, []firestore.Update{{Path: "fcm_tokens", Value: tokensMap}})
	if err != nil {
		log.Printf("Error saving fcm_tokens back to user doc: %v", err)
	}
	return nil
}
