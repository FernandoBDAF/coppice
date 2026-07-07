package profile

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/fernandobarroso/microservices/operational-workers/internal/common/queue"
)

var (
	ErrInvalidProfileID = errors.New("invalid profile_id")
	ErrInvalidTaskType  = errors.New("invalid task_type")
	ErrInvalidEnvelope  = errors.New("envelope type must be 'profile.task'")
)

// TaskType is advisory only: MESSAGE_FORMAT.md does not close this to an
// enum, so unrecognized values are handled generically rather than
// rejected (forward compatibility).
type TaskType string

const (
	TaskTypeSync     TaskType = "sync"
	TaskTypeValidate TaskType = "validate"
	TaskTypeEnrich   TaskType = "enrich"
)

// ProfilePayload mirrors graph-worker/shared/contracts/MESSAGE_FORMAT.md
// "Profile Payload" exactly.
type ProfilePayload struct {
	TaskType  string                 `json:"task_type"`
	ProfileID string                 `json:"profile_id"`
	UserID    string                 `json:"user_id,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// ProfileMessage is the parsed envelope + typed payload for a profile.task task.
type ProfileMessage struct {
	ID        string
	Type      string
	Timestamp time.Time
	Payload   ProfilePayload
	Metadata  map[string]string
}

// NewProfileMessage decodes msg.Payload (the envelope's "payload" object)
// into ProfilePayload.
func NewProfileMessage(msg *queue.Message) (*ProfileMessage, error) {
	var payload ProfilePayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return nil, err
	}

	return &ProfileMessage{
		ID:        msg.ID,
		Type:      msg.Type,
		Timestamp: msg.Timestamp,
		Payload:   payload,
		Metadata:  msg.Metadata,
	}, nil
}

// Validate checks the envelope + required payload fields. Unknown
// task_type values are tolerated (see TaskType doc); only structurally
// required fields are enforced.
func (m *ProfileMessage) Validate() error {
	if m.Type != "" && m.Type != "profile.task" {
		return ErrInvalidEnvelope
	}

	if m.Payload.ProfileID == "" {
		return ErrInvalidProfileID
	}

	if m.Payload.TaskType == "" {
		return ErrInvalidTaskType
	}

	return nil
}
