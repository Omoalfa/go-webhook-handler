package events

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Processor struct {
	repository Repository
}

func NewProcessor(repository Repository) *Processor {
	return &Processor{
		repository: repository,
	}
}

func (p *Processor) Process(
	ctx context.Context,
	id uuid.UUID,
) error {
	event, err := p.repository.Get(ctx, id)
	if err != nil {
		return err
	}

	if event.Status == EventStatusProcessed {
		return nil
	}

	_, err = p.repository.UpdateStatus(
		ctx,
		event.Id,
		string(EventStatusProcessing),
	)
	if err != nil {
		return err
	}

	select {
	case <-time.After(3 * time.Second):
	case <-ctx.Done():
		return ctx.Err()
	}

	_, err = p.repository.UpdateStatus(
		ctx,
		event.Id,
		string(EventStatusProcessed),
	)
	if err != nil {
		return err
	}

	return nil
}
