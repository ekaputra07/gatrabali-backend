package gatrabali

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"gatrabali/constant"
	"gatrabali/model"
	"gatrabali/sync"
)

// PubSubMessage is the payload of Pub/Sub message
type PubSubMessage struct {
	Data []byte `json:"data"`
}

func (m *PubSubMessage) getPayload() (*model.Payload, error) {
	var payload model.Payload
	if err := json.Unmarshal(m.Data, &payload); err != nil {
		return nil, err
	}
	if payload.ID == nil || payload.Type == nil || payload.Op == nil {
		return nil, errors.New("Invalid message payload (missing id, type or op)")
	}
	return &payload, nil
}

// SyncData calls Miniflux API and store its response to Cloud Firestore
func SyncData(ctx context.Context, m PubSubMessage) error {
	log.Printf("SyncData triggered with payload: %v\n", string(m.Data))

	payload, err := m.getPayload()
	if err != nil {
		return err
	}

	firestore, err := Firestore()
	if err != nil {
		return err
	}
	defer firestore.Close()

	switch *payload.Type {
	case constant.TypeCategory:
		return sync.StartCategorySync(firestore, payload)
	case constant.TypeFeed:
		return sync.StartFeedSync(firestore, payload)
	case constant.TypeEntry:
		return sync.StartEntrySync(firestore, payload)
	}
	return nil
}
