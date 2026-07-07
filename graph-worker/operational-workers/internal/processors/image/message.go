package image

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/fernandobarroso/microservices/operational-workers/internal/common/queue"
)

var (
	ErrInvalidSourceURL  = errors.New("invalid source_url")
	ErrInvalidOperation  = errors.New("invalid operation")
	ErrInvalidTargetPath = errors.New("invalid target_path")
	ErrInvalidDimensions = errors.New("width/height must be positive when set")
	ErrInvalidQuality    = errors.New("quality must be between 1 and 100 when set")
	ErrInvalidEnvelope   = errors.New("envelope type must be 'image.process'")
)

// Operation is advisory only: MESSAGE_FORMAT.md does not close this to an
// enum, so unrecognized values are simulated generically rather than
// rejected (forward compatibility).
type Operation string

const (
	OperationResize   Operation = "resize"
	OperationCompress Operation = "compress"
	OperationConvert  Operation = "convert"
)

// ImagePayload mirrors graph-worker/shared/contracts/MESSAGE_FORMAT.md
// "Image Payload" exactly.
type ImagePayload struct {
	Operation  string `json:"operation"`
	SourceURL  string `json:"source_url"`
	TargetPath string `json:"target_path"`
	Width      int    `json:"width,omitempty"`
	Height     int    `json:"height,omitempty"`
	Quality    int    `json:"quality,omitempty"`
	Format     string `json:"format,omitempty"`
}

// ImageMessage is the parsed envelope + typed payload for an image.process task.
type ImageMessage struct {
	ID        string
	Type      string
	Timestamp time.Time
	Payload   ImagePayload
	Metadata  map[string]string
}

// NewImageMessage decodes msg.Payload (the envelope's "payload" object)
// into ImagePayload.
func NewImageMessage(msg *queue.Message) (*ImageMessage, error) {
	var payload ImagePayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return nil, err
	}

	return &ImageMessage{
		ID:        msg.ID,
		Type:      msg.Type,
		Timestamp: msg.Timestamp,
		Payload:   payload,
		Metadata:  msg.Metadata,
	}, nil
}

// Validate checks the envelope + required payload fields. Unknown
// operation values are tolerated; only structurally required/well-formed
// fields are enforced.
func (m *ImageMessage) Validate() error {
	if m.Type != "" && m.Type != "image.process" {
		return ErrInvalidEnvelope
	}

	if m.Payload.SourceURL == "" {
		return ErrInvalidSourceURL
	}

	if m.Payload.Operation == "" {
		return ErrInvalidOperation
	}

	if m.Payload.TargetPath == "" {
		return ErrInvalidTargetPath
	}

	if m.Payload.Width < 0 || m.Payload.Height < 0 {
		return ErrInvalidDimensions
	}

	if m.Payload.Quality != 0 && (m.Payload.Quality < 1 || m.Payload.Quality > 100) {
		return ErrInvalidQuality
	}

	return nil
}
