package events

import (
	"context"
	"errors"
	"time"

	"github.com/Omoalfa/event-platform/internal/queue"
	"github.com/google/uuid"
)

var (
	ErrEventNotFound = errors.New("Event not found")
	ErrInvalidStatus = errors.New("Invalid event status")
)

type Service struct {
	repository Repository
	queue      queue.Queue
}

func NewService(
	repository Repository,
	queue queue.Queue,
) *Service {
	return &Service{
		repository: repository,
		queue:      queue,
	}
}

func (s *Service) Create(ctx context.Context, cEvent CreateEvent) (Event, error) {
	if cEvent.Provider == "" {
		return Event{}, errors.New("provider is required")
	}

	if cEvent.EventType == "" {
		return Event{}, errors.New("event type is required")
	}

	event, err := s.repository.Create(ctx, cEvent)
	if err != nil {
		return Event{}, err
	}

	if err := s.queue.EnqueueEvent(ctx, event.Id); err != nil {
		return Event{}, err
	}

	return event, nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (Event, error) {
	return s.repository.Get(ctx, id)
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]Event, error) {
	if limit <= 0 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	return s.repository.List(ctx, limit, offset)
}

func (s *Service) UpdateStatus(
	ctx context.Context,
	id uuid.UUID,
	status string,
) (Event, error) {
	if !isValidStatus(status) {
		return Event{}, ErrInvalidStatus
	}

	return s.repository.UpdateStatus(ctx, id, status)
}

func isValidStatus(status string) bool {
	switch status {
	case "pending", "processing", "processed", "failed":
		return true
	default:
		return false
	}
}

func (s *Service) Process(
	ctx context.Context,
	id uuid.UUID,
) error {
	event, err := s.repository.Get(ctx, id)
	if err != nil {
		return err
	}

	_, err = s.repository.UpdateStatus(
		ctx,
		event.Id,
		string(EventStatusProcessing),
	)
	if err != nil {
		return err
	}

	//simulate task processing::::
	select {
	case <-time.After(3 * time.Second):
	case <-ctx.Done():
		return ctx.Err()
	}

	_, err = s.repository.UpdateStatus(
		ctx,
		event.Id,
		string(EventStatusProcessed),
	)
	if err != nil {
		return err
	}

	return nil
}
