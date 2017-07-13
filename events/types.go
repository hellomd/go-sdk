package events

var ContextKey = struct{}{}

// Publisher is a client that can publish events
type Publisher interface {
	// Publish publishes an event
	Publish(key string, body interface{}) error
}
