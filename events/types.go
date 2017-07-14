package events

// ContextKey is the key that should be used for storing a client in the context
var ContextKey = struct{}{}

// Client can publish and subscribe to events
type Client interface {
	Publisher
	Subscriber
}

// Publisher is a client that can publish events
type Publisher interface {
	// Publish publishes an event
	Publish(key string, body interface{}) error
}

// Subscriber is a client that can subscribe to events
type Subscriber interface {
	// Subscribe
	Subscribe(pattern string) (Subscription, error)
}

// Subscription can receive events to which it subscribes
type Subscription interface {
	// Receive returns a channel to which subscription events are passed
	Receive() <-chan []byte

	// NotifyClose notifies of when the subscription is ended
	NotifyClose() <-chan struct{}

	// Close stops receiving events for this subscription
	Close()

	// Errors returns a channel where subscription errors are logged
	Errors() <-chan error
}
