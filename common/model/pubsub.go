package model

// SyncPayload is the pub/sub message payload after deserialized
type SyncPayload struct {
	ID   *int64  `json:"entity_id"`
	Type *string `json:"entity_type"`
	Op   *string `json:"entity_op"`
}

// PushNotificationPayload is the payload to send push notification.
// Title and Meta can be omitted but Message is mandatory.
type PushNotificationPayload struct {
	Title   *string
	Message *string
	Meta    *map[string]interface{}
}
