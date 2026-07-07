package email

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/fernandobarroso/microservices/operational-workers/internal/common/queue"
)

var (
	ErrInvalidRecipient = errors.New("invalid recipient email")
	ErrInvalidEmailType = errors.New("invalid email type")
	ErrInvalidEnvelope  = errors.New("envelope type must be 'email.send'")
)

// EmailType is advisory only: MESSAGE_FORMAT.md does not close this to an
// enum, so unrecognized values are simulated generically rather than
// rejected (forward compatibility).
type EmailType string

const (
	EmailTypeWelcome      EmailType = "welcome"
	EmailTypeNotification EmailType = "notification"
	EmailTypeAlert        EmailType = "alert"
)

// EmailPayload mirrors graph-worker/shared/contracts/MESSAGE_FORMAT.md
// "Email Payload" exactly.
type EmailPayload struct {
	EmailType  string                 `json:"email_type"`
	Recipient  string                 `json:"recipient"`
	Subject    string                 `json:"subject,omitempty"`
	TemplateID string                 `json:"template_id,omitempty"`
	Variables  map[string]interface{} `json:"variables,omitempty"`
}

// EmailMessage is the parsed envelope + typed payload for an email.send task.
type EmailMessage struct {
	ID        string
	Type      string
	Timestamp time.Time
	Payload   EmailPayload
	Metadata  map[string]string
}

// NewEmailMessage decodes msg.Payload (the envelope's "payload" object,
// already extracted by the queue layer) into EmailPayload. It does NOT
// re-unmarshal the whole envelope — msg.Payload is just the inner object,
// per graph-worker/shared/contracts/MESSAGE_FORMAT.md.
func NewEmailMessage(msg *queue.Message) (*EmailMessage, error) {
	var payload EmailPayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return nil, err
	}

	return &EmailMessage{
		ID:        msg.ID,
		Type:      msg.Type,
		Timestamp: msg.Timestamp,
		Payload:   payload,
		Metadata:  msg.Metadata,
	}, nil
}

// Validate checks the envelope + required payload fields. Unknown
// email_type values are tolerated (see EmailType doc); only structurally
// required fields are enforced.
func (m *EmailMessage) Validate() error {
	if m.Type != "" && m.Type != "email.send" {
		return ErrInvalidEnvelope
	}

	if m.Payload.Recipient == "" || !strings.Contains(m.Payload.Recipient, "@") {
		return ErrInvalidRecipient
	}

	if m.Payload.EmailType == "" {
		return ErrInvalidEmailType
	}

	return nil
}
