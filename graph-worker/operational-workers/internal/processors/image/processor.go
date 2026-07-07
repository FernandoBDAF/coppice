package image

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fernandobarroso/microservices/operational-workers/internal/common/queue"
	"github.com/fernandobarroso/microservices/operational-workers/internal/common/utils"
)

// ImageProcessor handles image processing messages. No real image
// manipulation is performed (per mission scope) — Process simulates the
// operation and logs coherently so behavior is observable end-to-end.
type ImageProcessor struct {
	metrics *utils.ProcessorMetrics
}

// NewImageProcessor creates a new image processor
func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{
		metrics: utils.NewProcessorMetrics("image"),
	}
}

// Process processes an image message
func (p *ImageProcessor) Process(ctx context.Context, msg *queue.Message) error {
	timer := p.metrics.StartTimer()
	defer timer.ObserveDuration()
	p.metrics.RecordProcessingStart()

	imageMsg, err := NewImageMessage(msg)
	if err != nil {
		p.metrics.RecordProcessingError()
		return fmt.Errorf("failed to parse image message: %w", err)
	}

	if err := imageMsg.Validate(); err != nil {
		p.metrics.RecordProcessingError()
		return err
	}

	if err := p.runOperation(ctx, imageMsg); err != nil {
		p.metrics.RecordProcessingError()
		return err
	}

	p.metrics.RecordProcessingSuccess()
	return nil
}

// Validate validates the message
func (p *ImageProcessor) Validate(msg *queue.Message) error {
	imageMsg, err := NewImageMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to parse message for validation: %w", err)
	}
	return imageMsg.Validate()
}

// Type returns the processor type
func (p *ImageProcessor) Type() string {
	return "image"
}

// HandleError handles processing errors. Returning a non-nil error here
// tells the consumer to nack the message WITHOUT requeue, so it lands on
// the image-processing.dlq instead of being redelivered forever.
func (p *ImageProcessor) HandleError(ctx context.Context, msg *queue.Message, err error) error {
	log.Printf("image processing error: %v", err)
	return err
}

// runOperation simulates the requested image operation. Unrecognized
// operation values still get a generic simulated pass (forward
// compatible) rather than failing.
func (p *ImageProcessor) runOperation(ctx context.Context, msg *ImageMessage) error {
	delay := 500 * time.Millisecond

	switch Operation(msg.Payload.Operation) {
	case OperationResize:
		log.Printf("resizing %s -> %s (%dx%d)", msg.Payload.SourceURL, msg.Payload.TargetPath, msg.Payload.Width, msg.Payload.Height)
	case OperationCompress:
		log.Printf("compressing %s -> %s (quality=%d)", msg.Payload.SourceURL, msg.Payload.TargetPath, msg.Payload.Quality)
	case OperationConvert:
		log.Printf("converting %s -> %s (format=%s)", msg.Payload.SourceURL, msg.Payload.TargetPath, msg.Payload.Format)
	default:
		log.Printf("processing image operation=%s %s -> %s", msg.Payload.Operation, msg.Payload.SourceURL, msg.Payload.TargetPath)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(delay):
		log.Printf("image operation %q completed for %s", msg.Payload.Operation, msg.Payload.SourceURL)
		return nil
	}
}
