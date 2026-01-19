package broker

import (
	"context"

	"github.com/google/uuid"
)

type Publisher interface {

	//TODO: uprgrade message to Avro - for efficiency
	Publish(context.Context, string, string, []byte) error
}

type PublisherPayload struct {
	EventID uuid.UUID   `json:"event_id"`
	Payload interface{} `json:"payload"`
}
