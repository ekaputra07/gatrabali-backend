package function

import (
	"context"
	"encoding/json"
	"errors"
	"log"

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

	var payload model.Payload
	if err := json.Unmarshal(m.Data, &payload); err != nil {
		return err
	}
	if payload.ID == nil || payload.Type == nil || payload.Op == nil {
		return errors.New("Invalid message payload (missing id, type or op)")
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
	return nil
}
