package email

import (
	"encoding/json"
	"testing"

	"github.com/fernandobarroso/microservices/operational-workers/internal/common/queue"
)

func TestNewEmailMessage_HappyPath(t *testing.T) {
	envelope := []byte(`{
		"id": "11111111-1111-1111-1111-111111111111",
		"type": "email.send",
		"timestamp": "2026-01-30T12:34:56Z",
		"payload": {
			"email_type": "welcome",
			"recipient": "user@example.com",
			"subject": "Welcome",
			"template_id": "welcome-template",
			"variables": {"first_name": "Ada"}
		}
	}`)

	var msg queue.Message
	if err := json.Unmarshal(envelope, &msg); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}

	emailMsg, err := NewEmailMessage(&msg)
	if err != nil {
		t.Fatalf("NewEmailMessage: %v", err)
	}

	if emailMsg.Payload.Recipient != "user@example.com" {
		t.Errorf("Recipient = %q, want user@example.com", emailMsg.Payload.Recipient)
	}
	if emailMsg.Payload.EmailType != "welcome" {
		t.Errorf("EmailType = %q, want welcome", emailMsg.Payload.EmailType)
	}
	if emailMsg.Payload.Variables["first_name"] != "Ada" {
		t.Errorf("Variables[first_name] = %v, want Ada", emailMsg.Payload.Variables["first_name"])
	}

	if err := emailMsg.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

func TestEmailMessage_Validate_MissingRecipient(t *testing.T) {
	msg := &EmailMessage{Type: "email.send", Payload: EmailPayload{EmailType: "welcome"}}
	if err := msg.Validate(); err == nil {
		t.Error("expected error for missing recipient")
	}
}

func TestEmailMessage_Validate_WrongEnvelopeType(t *testing.T) {
	msg := &EmailMessage{
		Type:    "image.process",
		Payload: EmailPayload{Recipient: "a@b.com", EmailType: "welcome"},
	}
	if err := msg.Validate(); err == nil {
		t.Error("expected error for mismatched envelope type")
	}
}

func TestEmailMessage_Validate_UnknownEmailTypeTolerated(t *testing.T) {
	// email_type is not a closed enum in MESSAGE_FORMAT.md; unrecognized
	// values must not be treated as poison.
	msg := &EmailMessage{
		Type:    "email.send",
		Payload: EmailPayload{Recipient: "a@b.com", EmailType: "future-type"},
	}
	if err := msg.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil for forward-compatible email_type", err)
	}
}
