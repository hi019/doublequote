package asynq

import (
	"context"
	"encoding/json"
	"time"

	dq "doublequote/pkg/domain"
	"github.com/hibiken/asynq"
)

// Ensure type implements interface.
var _ dq.EventService = (*EventService)(nil)

type EventService struct {
	server *asynq.Server
	mux    *asynq.ServeMux
	client *asynq.Client
}

func NewEventService(redisUrl string) *EventService {
	s := EventService{
		mux: asynq.NewServeMux(),
	}

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisUrl},
		asynq.Config{Concurrency: 10},
	)
	s.server = srv

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisUrl})
	s.client = client

	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: redisUrl},
		&asynq.SchedulerOpts{
			Location: time.UTC,
		},
	)

	scheduler.Register("a")

	return &s
}

func (s *EventService) Subscribe(topic string, handler dq.EventHandler) {
	s.mux.HandleFunc(topic, func(ctx context.Context, task *asynq.Task) error {
		return handler(dq.Event{
			Topic:   topic,
			Payload: task.Payload(),
		})
	})
}

func (s *EventService) Publish(topic string, payload dq.Payload) error {
	newPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = s.client.Enqueue(asynq.NewTask(topic, newPayload))
	return err
}

// Open starts the Asynq server.
// NOTE: it must be called after all event handlers are subscribed.
func (s *EventService) Open() error {
	return s.server.Start(s.mux)
}

func (s *EventService) Close() error {
	s.server.Shutdown()
	return s.client.Close()
}
