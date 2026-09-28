package queue

import (
	"context"

	"github.com/google/uuid"
)

type Queue interface {
	EnqueueEvent(ctx context.Context, eventID uuid.UUID) error
}
