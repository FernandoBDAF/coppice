package task

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Publisher interface {
	PublishWithRoutingKey(routingKey string, msg *Message) error
}

type Service struct {
	publisher Publisher
}

func NewService(publisher Publisher) *Service {
	return &Service{publisher: publisher}
}

func (s *Service) Submit(ctx context.Context, routingKey, msgType string, payload interface{}, metadata map[string]string) (string, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	correlationID := uuid.New().String()

	// Every published message carries metadata.source="api-service" and a
	// trace_id, per the pinned envelope shape. Callers may supply their own
	// trace_id (e.g. propagated from an inbound request); otherwise fall back
	// to this publish's correlation ID so messages stay traceable end to end.
	if metadata == nil {
		metadata = map[string]string{}
	}
	metadata["source"] = "api-service"
	if _, ok := metadata["trace_id"]; !ok {
		metadata["trace_id"] = correlationID
	}

	msg := &Message{
		ID:            uuid.New().String(),
		Type:          msgType,
		Timestamp:     time.Now().UTC(),
		CorrelationID: correlationID,
		Payload:       body,
		Metadata:      metadata,
		Priority:      0,
	}

	if err := s.publisher.PublishWithRoutingKey(routingKey, msg); err != nil {
		return "", err
	}

	return msg.ID, nil
}
