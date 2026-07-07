package task

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

type mockPublisher struct {
	lastRoutingKey string
	lastMessage    *Message
	err            error
	calls          int
}

func (m *mockPublisher) PublishWithRoutingKey(routingKey string, msg *Message) error {
	m.calls++
	m.lastRoutingKey = routingKey
	m.lastMessage = msg
	return m.err
}

func TestService_Submit_BuildsEnvelopeAndPublishes(t *testing.T) {
	pub := &mockPublisher{}
	svc := NewService(pub)

	payload := map[string]interface{}{"profile_id": "abc-123"}
	metadata := map[string]string{"source": "api-service"}

	taskID, err := svc.Submit(context.Background(), "profile.task", "profile.task", payload, metadata)
	if err != nil {
		t.Fatalf("Submit returned error: %v", err)
	}
	if taskID == "" {
		t.Fatalf("expected non-empty task ID")
	}
	if pub.calls != 1 {
		t.Fatalf("expected exactly one publish call, got %d", pub.calls)
	}
	if pub.lastRoutingKey != "profile.task" {
		t.Errorf("expected routing key 'profile.task', got %q", pub.lastRoutingKey)
	}
	if pub.lastMessage.ID != taskID {
		t.Errorf("expected message ID to match returned task ID")
	}
	if pub.lastMessage.Type != "profile.task" {
		t.Errorf("expected message type 'profile.task', got %q", pub.lastMessage.Type)
	}
	if pub.lastMessage.Timestamp.IsZero() {
		t.Errorf("expected message timestamp to be set")
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(pub.lastMessage.Payload, &decoded); err != nil {
		t.Fatalf("failed to decode published payload: %v", err)
	}
	if decoded["profile_id"] != "abc-123" {
		t.Errorf("expected payload to round-trip profile_id, got %+v", decoded)
	}
}

func TestService_Submit_PropagatesPublishError(t *testing.T) {
	wantErr := errors.New("broker unavailable")
	pub := &mockPublisher{err: wantErr}
	svc := NewService(pub)

	taskID, err := svc.Submit(context.Background(), "email.send", "email.send", map[string]interface{}{}, nil)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected publish error to propagate, got %v", err)
	}
	if taskID != "" {
		t.Errorf("expected empty task ID on failure, got %q", taskID)
	}
}
