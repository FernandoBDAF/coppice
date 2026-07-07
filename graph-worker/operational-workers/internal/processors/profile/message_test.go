package profile

import (
	"encoding/json"
	"testing"

	"github.com/fernandobarroso/microservices/operational-workers/internal/common/queue"
)

func TestNewProfileMessage_HappyPath(t *testing.T) {
	envelope := []byte(`{
		"id": "id-1",
		"type": "profile.task",
		"timestamp": "2026-01-30T12:34:56Z",
		"payload": {
			"task_type": "sync",
			"profile_id": "profile-789",
			"user_id": "user-456",
			"data": {"source": "external-system"}
		}
	}`)

	var msg queue.Message
	if err := json.Unmarshal(envelope, &msg); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}

	profileMsg, err := NewProfileMessage(&msg)
	if err != nil {
		t.Fatalf("NewProfileMessage: %v", err)
	}

	if profileMsg.Payload.ProfileID != "profile-789" || profileMsg.Payload.TaskType != "sync" {
		t.Errorf("payload = %+v", profileMsg.Payload)
	}

	if err := profileMsg.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

func TestProfileMessage_Validate_MissingProfileID(t *testing.T) {
	msg := &ProfileMessage{Type: "profile.task", Payload: ProfilePayload{TaskType: "sync"}}
	if err := msg.Validate(); err == nil {
		t.Error("expected error for missing profile_id")
	}
}

func TestProfileMessage_Validate_UnknownTaskTypeTolerated(t *testing.T) {
	// task_type is not a closed enum in MESSAGE_FORMAT.md; unrecognized
	// values must not be treated as poison.
	msg := &ProfileMessage{Type: "profile.task", Payload: ProfilePayload{ProfileID: "p-1", TaskType: "future-task"}}
	if err := msg.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil for forward-compatible task_type", err)
	}
}
