package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"

	"github.com/Omoalfa/event-platform/internal/events"
	"github.com/Omoalfa/event-platform/internal/queue"
)

type EventProcessor struct {
	processor *events.Processor
}

func NewEventProcessor(processor *events.Processor) *EventProcessor {
	return &EventProcessor{
		processor: processor,
	}
}

func (p *EventProcessor) ProcessTask(
	ctx context.Context,
	task *asynq.Task,
) error {
	var payload queue.EventTaskPayload

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("decode event task payload: %w", err)
	}

	eventID, err := uuid.Parse(payload.EventID)
	if err != nil {
		return fmt.Errorf("parse event ID: %w", err)
	}

	if err := p.processor.Process(ctx, eventID); err != nil {
		return fmt.Errorf("process event: %w", err)
	}

	fmt.Printf("Processed event %s\n", eventID)

	return nil
}
