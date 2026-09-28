package events

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, event CreateEvent) (Event, error)
	Get(ctx context.Context, id uuid.UUID) (Event, error)
	List(ctx context.Context, limit, offset int) ([]Event, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (Event, error)
}
