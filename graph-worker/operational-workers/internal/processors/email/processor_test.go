package email

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/fernandobarroso/microservices/operational-workers/internal/common/queue"
)

func TestEmailProcessor_Process_HappyPath(t *testing.T) {
	payload, err := json.Marshal(EmailPayload{
		EmailType:  "welcome",
		Recipient:  "user@example.com",
		TemplateID: "welcome-template",
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	msg := &queue.Message{
		ID:      "11111111-1111-1111-1111-111111111111",
		Type:    "email.send",
		Payload: payload,
	}

	p := NewEmailProcessor()

	if err := p.Validate(msg); err != nil {
		t.Fatalf("Validate() error = %v, want nil", err)
	}

	if err := p.Process(context.Background(), msg); err != nil {
		t.Fatalf("Process() error = %v, want nil", err)
	}
}

func TestEmailProcessor_Process_PoisonMessageErrors(t *testing.T) {
	msg := &queue.Message{
		ID:      "id-2",
		Type:    "email.send",
		Payload: json.RawMessage(`{"email_type":"welcome"}`), // missing recipient
	}

	p := NewEmailProcessor()

	if err := p.Process(context.Background(), msg); err == nil {
		t.Error("expected Process to fail for a payload missing the required recipient field")
	}
}
