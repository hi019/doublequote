package asynq

import (
	"context"
	"encoding/json"

	domain "doublequote/pkg/domain"
	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

// Ensure type implements interface
var _ domain.EventService = (*EventService)(nil)

type EventService struct {
	server    *asynq.Server
	mux       *asynq.ServeMux
	client    *asynq.Client
	scheduler *asynq.Scheduler
}

func NewEventService(lc fx.Lifecycle, cfg domain.Config) *EventService {
	s := EventService{
		mux: asynq.NewServeMux(),
	}

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.Redis.URL},
		asynq.Config{Concurrency: 10},
	)
	s.server = srv

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.Redis.URL})
	s.client = client

	scheduler := asynq.NewScheduler(asynq.RedisClientOpt{Addr: cfg.Redis.URL}, nil)
	s.scheduler = scheduler

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return s.open()
		},
		OnStop: func(ctx context.Context) error {
			return s.close()
		},
	})

	return &s
}

func (s *EventService) Subscribe(topic string, handler domain.EventHandler) {
	s.mux.HandleFunc(topic, func(ctx context.Context, task *asynq.Task) error {
		return handler(domain.Event{
			Topic:   topic,
			Payload: task.Payload(),
		})
	})
}

func (s *EventService) Publish(topic string, payload domain.Payload) error {
	newPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = s.client.Enqueue(asynq.NewTask(topic, newPayload))
	return err
}

func (s *EventService) PublishPeriodic(topic, cron string, payload domain.Payload) error {
	// TODO test

	newPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if _, err = s.scheduler.Register(cron, asynq.NewTask(topic, newPayload)); err != nil {
		return err
	}

	return nil
}

// open starts the Asynq server
// NOTE: it must be called after all event handlers are subscribed
func (s *EventService) open() error {
	return s.server.Start(s.mux)
}

func (s *EventService) close() error {
	s.server.Shutdown()
	return s.client.Close()
}
