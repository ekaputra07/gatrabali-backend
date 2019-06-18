package model

// Payload is the pub/sub message payload after deserialized
type Payload struct {
	ID   *int64  `json:"entity_id"`
	Type *string `json:"entity_type"`
	Op   *string `json:"entity_op"`
}
