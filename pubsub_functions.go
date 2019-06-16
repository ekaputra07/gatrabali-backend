package gatrabali

import (
	// "log"
	"context"
)

// PubSubMessage is the payload of Pub/Sub message
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// SyncData calls Miniflux API and store its response to Cloud Firestore
func SyncData(ctx context.Context, m PubSubMessage) error {
	return nil
}
