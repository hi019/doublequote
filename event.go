package dq

// Event topic constants.
const (
	EventTopicUserCreated = "user:created"
)

type EventHandler func(event Event) error

type Payload interface{}

// Event represents an event that occurs in the system.
type Event struct {
	// Specifies the type of event that is occurring.
	Topic string

	// The actual data from the event. See related payload types below.
	Payload []byte
}

// UserCreatedPayload represents the payload for an Event object with a
// type of EventTopicUserCreated.
type UserCreatedPayload struct {
	User *User
}

// EventService represents a service for managing event dispatch and event
// listeners (aka subscriptions).
type EventService interface {
	// Publish publishes an event to a user's event listeners.
	// If the user is not currently subscribed then this is a no-op.
	Publish(topic string, payload Payload) error

	// Subscribe creates a subscription for the current user's events.
	// Caller must call Subscription.Close() when done with the subscription.
	Subscribe(topic string, handler EventHandler)
}
