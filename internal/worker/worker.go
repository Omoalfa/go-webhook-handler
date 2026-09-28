package worker

import (
	"github.com/Omoalfa/event-platform/internal/queue"
	"github.com/hibiken/asynq"
)

type Worker struct {
	server *asynq.Server
}

func NewWorker(redisURL string) (*Worker, error) {
	redisOpts, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		return nil, err
	}

	server := asynq.NewServer(
		redisOpts,
		asynq.Config{
			Concurrency: 10,
		},
	)

	return &Worker{
		server: server,
	}, nil
}

func (w *Worker) Start(processor *EventProcessor) error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(
		queue.TaskProcessEvent,
		processor.ProcessTask,
	)

	return w.server.Run(mux)
}

func (w *Worker) Shutdown() {
	w.server.Shutdown()
}
