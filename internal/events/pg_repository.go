package events

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Omoalfa/event-platform/internal/database/db"
)

type PGRepository struct {
	queries *db.Queries
}

func NewPGRepository(pool *pgxpool.Pool) *PGRepository {
	return &PGRepository{
		queries: db.New(pool),
	}
}

func (r *PGRepository) Create(ctx context.Context, c_event CreateEvent) (Event, error) {
	row, err := r.queries.CreateEvent(ctx, db.CreateEventParams{
		Provider:  c_event.Provider,
		EventType: c_event.EventType,
		Payload:   c_event.Payload,
		Status:    "pending",
	})
	if err != nil {
		return Event{}, err
	}

	event, err := toDomainEvent(row)
	if err != nil {
		return Event{}, err
	}

	return event, nil
}

func (r *PGRepository) Get(ctx context.Context, id uuid.UUID) (Event, error) {
	row, err := r.queries.GetEvent(ctx, pgtype.UUID{
		Bytes: id,
		Valid: true,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Event{}, ErrEventNotFound
		}
		return Event{}, err
	}

	event, err := toDomainEvent(row)
	if err != nil {
		return Event{}, err
	}

	return event, nil
}

func (r *PGRepository) List(ctx context.Context, limit, offset int) ([]Event, error) {
	rows, err := r.queries.ListEvents(ctx, db.ListEventsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	events := make([]Event, 0, len(rows))

	for _, row := range rows {
		event, err := toDomainEvent(row)
		if err != nil {
			return events, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (r *PGRepository) UpdateStatus(
	ctx context.Context,
	id uuid.UUID,
	status string,
) (Event, error) {
	row, err := r.queries.UpdateEventStatus(ctx, db.UpdateEventStatusParams{
		ID: pgtype.UUID{
			Bytes: id,
			Valid: true,
		},
		Status: status,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Event{}, ErrEventNotFound
		}
		return Event{}, err
	}

	event, err := toDomainEvent(row)
	if err != nil {
		return Event{}, err
	}

	return event, nil
}

func toDomainEvent(row db.Event) (Event, error) {
	var payload map[string]interface{}
	err := json.Unmarshal(row.Payload, &payload)
	if err != nil {
		return Event{}, err
	}
	return Event{
		Id:        uuid.UUID(row.ID.Bytes),
		Provider:  row.Provider,
		EventType: row.EventType,
		Payload:   &payload,
		Status:    toEventStatus(row.Status),
		CreatedAt: &row.CreatedAt.Time,
		UpdatedAt: &row.UpdatedAt.Time,
	}, nil
}

func toEventStatus(status string) EventStatus {
	switch status {
	case "failed":
		return EventStatusFailed
	case "pending":
		return EventStatusPending
	case "processed":
		return EventStatusProcessed
	case "processing":
		return EventStatusProcessing
	default:
		return ""
	}
}

// func toDbEvent(event Event) (db.Event, error) {
// 	payload, err := json.Marshal(event.Payload)
// 	if err != nil {
// 		return db.Event{}, err
// 	}

// 	return db.Event{
// 		ID:        pgtype.UUID{Bytes: event.Id, Valid: true},
// 		Provider:  event.Provider,
// 		EventType: event.EventType,
// 		Payload:   payload,
// 		Status:    string(event.Status),
// 		CreatedAt: pgtype.Timestamptz{Time: *event.CreatedAt, Valid: true},
// 		UpdatedAt: pgtype.Timestamptz{Time: *event.UpdatedAt, Valid: true},
// 	}, nil
// }
