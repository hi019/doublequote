package asynq

import (
	"os"
	"testing"
	"time"

	dq "doublequote/pkg/domain"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO I think this is more of an integration/feature test. Should it be put somewhere else?

func TestEventService(t *testing.T) {
	err := godotenv.Load("../.env.testing")
	require.Nil(t, err)

	t.Run("OK", func(t *testing.T) {
		es := NewEventService(os.Getenv("REDIS_URL"))
		defer func(es *EventService) {
			assert.Nil(t, es.Close())
		}(es)

		err := es.Open()
		assert.Nil(t, err)

		// Wait until server is ready
		time.Sleep(3 * time.Second)

		var event dq.Event
		es.Subscribe("test", func(e dq.Event) error {
			event = e
			return nil
		})

		err = es.Publish("test", "test")
		assert.Nil(t, err)

		// Give the event time to process
		time.Sleep(2 * time.Second)

		assert.Equal(t, "test", event.Topic)
	})
}
