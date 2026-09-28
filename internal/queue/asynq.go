package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type AsynqQueue struct {
	client *asynq.Client
}

func NewAsynqQueue(redisURL string) (*AsynqQueue, error) {
	redisOpts, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	client := asynq.NewClient(redisOpts)

	return &AsynqQueue{
		client: client,
	}, nil
}

func (q *AsynqQueue) EnqueueEvent(
	ctx context.Context,
	eventID uuid.UUID,
) error {
	payload, err := json.Marshal(EventTaskPayload{
		EventID: eventID.String(),
	})
	if err != nil {
		return fmt.Errorf("marshal event task payload: %w", err)
	}

	task := asynq.NewTask(
		TaskProcessEvent,
		payload,
		asynq.MaxRetry(5),
	)

	_, err = q.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("enqueue event task: %w", err)
	}

	return nil
}

func (q *AsynqQueue) Close() error {
	return q.client.Close()
}
