package types

// SyncPayload is the pub/sub message payload after deserialized
type SyncPayload struct {
	ID   *int64  `json:"entity_id"`
	Type *string `json:"entity_type"`
	Op   *string `json:"entity_op"`
}

// PushNotificationPayload is the payload to send push notification.
type PushNotificationPayload struct {
	UserID      string            `json:"user_id"`
	Title       string            `json:"title"`
	Body        string            `json:"body"`
	CollapseKey string            `json:"collapse_key,omitempty"`
	Data        map[string]string `json:"data,omitempty"`
}
