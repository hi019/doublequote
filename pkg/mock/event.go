package mock

import (
	dq "doublequote/pkg/domain"
	"github.com/stretchr/testify/mock"
)

// Ensure type implements interface.
var _ dq.EventService = (*EventService)(nil)

type EventService struct {
	mock.Mock
}

func (s *EventService) Publish(topic string, payload dq.Payload) error {
	s.Called(topic, payload)

	return nil
}

func (s *EventService) Subscribe(topic string, handler dq.EventHandler) {
	s.Called(topic, handler)
}
