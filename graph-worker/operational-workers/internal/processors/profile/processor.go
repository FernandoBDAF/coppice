package profile

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fernandobarroso/microservices/operational-workers/internal/common/queue"
	"github.com/fernandobarroso/microservices/operational-workers/internal/common/utils"
)

// Processor handles profile task processing. No real profile store is
// wired up (per mission scope) — Process simulates the task and logs
// coherently so behavior is observable end-to-end.
type Processor struct {
	metrics *utils.ProcessorMetrics
}

// NewProcessor creates a new profile processor
func NewProcessor() *Processor {
	return &Processor{
		metrics: utils.NewProcessorMetrics("profile"),
	}
}

// Process processes a profile task message
func (p *Processor) Process(ctx context.Context, msg *queue.Message) error {
	timer := p.metrics.StartTimer()
	defer timer.ObserveDuration()
	p.metrics.RecordProcessingStart()

	profileMsg, err := NewProfileMessage(msg)
	if err != nil {
		p.metrics.RecordProcessingError()
		return fmt.Errorf("failed to parse profile message: %w", err)
	}

	if err := profileMsg.Validate(); err != nil {
		p.metrics.RecordProcessingError()
		return err
	}

	if err := p.runTask(ctx, profileMsg); err != nil {
		p.metrics.RecordProcessingError()
		return err
	}

	p.metrics.RecordProcessingSuccess()
	return nil
}

// Validate validates the message
func (p *Processor) Validate(msg *queue.Message) error {
	profileMsg, err := NewProfileMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to parse message for validation: %w", err)
	}
	return profileMsg.Validate()
}

// Type returns the processor type
func (p *Processor) Type() string {
	return "profile"
}

// HandleError handles processing errors. Returning a non-nil error here
// tells the consumer to nack the message WITHOUT requeue, so it lands on
// the profile-processing.dlq instead of being redelivered forever.
func (p *Processor) HandleError(ctx context.Context, msg *queue.Message, err error) error {
	log.Printf("profile processing error: %v", err)
	return err
}

// runTask dispatches to a task-specific simulated handler. Unrecognized
// task_type values still get a generic simulated pass (forward
// compatible) rather than failing.
func (p *Processor) runTask(ctx context.Context, msg *ProfileMessage) error {
	switch TaskType(msg.Payload.TaskType) {
	case TaskTypeSync:
		return p.simulate(ctx, msg, "syncing", 200*time.Millisecond, "synced")
	case TaskTypeValidate:
		return p.simulate(ctx, msg, "validating", 100*time.Millisecond, "validated")
	case TaskTypeEnrich:
		return p.simulate(ctx, msg, "enriching", 300*time.Millisecond, "enriched")
	default:
		log.Printf("unrecognized profile task_type %q for profile %s; handling generically", msg.Payload.TaskType, msg.Payload.ProfileID)
		return p.simulate(ctx, msg, "processing", 100*time.Millisecond, "processed")
	}
}

func (p *Processor) simulate(ctx context.Context, msg *ProfileMessage, verb string, delay time.Duration, doneVerb string) error {
	log.Printf("%s profile %s", verb, msg.Payload.ProfileID)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(delay):
		log.Printf("profile %s: %s", doneVerb, msg.Payload.ProfileID)
		return nil
	}
}
