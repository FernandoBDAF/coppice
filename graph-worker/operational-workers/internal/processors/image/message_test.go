package image

import (
	"encoding/json"
	"testing"

	"github.com/fernandobarroso/microservices/operational-workers/internal/common/queue"
)

func TestNewImageMessage_HappyPath(t *testing.T) {
	envelope := []byte(`{
		"id": "id-1",
		"type": "image.process",
		"timestamp": "2026-01-30T12:34:56Z",
		"payload": {
			"operation": "resize",
			"source_url": "s3://bucket/path/image.png",
			"target_path": "processed/image.png",
			"width": 512,
			"height": 512,
			"quality": 85,
			"format": "png"
		}
	}`)

	var msg queue.Message
	if err := json.Unmarshal(envelope, &msg); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}

	imageMsg, err := NewImageMessage(&msg)
	if err != nil {
		t.Fatalf("NewImageMessage: %v", err)
	}

	if imageMsg.Payload.Operation != "resize" || imageMsg.Payload.Width != 512 {
		t.Errorf("payload = %+v", imageMsg.Payload)
	}

	if err := imageMsg.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

func TestImageMessage_Validate_MissingSourceURL(t *testing.T) {
	msg := &ImageMessage{
		Type:    "image.process",
		Payload: ImagePayload{Operation: "resize", TargetPath: "x"},
	}
	if err := msg.Validate(); err == nil {
		t.Error("expected error for missing source_url")
	}
}

func TestImageMessage_Validate_QualityOutOfRange(t *testing.T) {
	msg := &ImageMessage{
		Type:    "image.process",
		Payload: ImagePayload{Operation: "compress", SourceURL: "s3://x", TargetPath: "y", Quality: 150},
	}
	if err := msg.Validate(); err == nil {
		t.Error("expected error for out-of-range quality")
	}
}
