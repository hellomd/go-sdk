package events

const (
	baseRetryDelay float64 = 10
	maxRetryDelay  float64 = 5000
)

// NewSubscriber creates a new client that can subscribe to events
//
// This subscriber is concurrent safe and has no initialization logic.
// Feel free to use whatever lifecycle you think is best for it.
func NewSubscriber(name, amqpURL string) *Subscriber {
	return &Subscriber{name, amqpURL}
}

// Subscriber is a client that can subscribe to events
type Subscriber struct {
	name    string
	amqpURL string
}

// Subscribe creates a new subscription for receiving events with keys that match the given pattern
func (s *Subscriber) Subscribe(pattern string) (*Subscription, error) {
	sub := &Subscription{
		subscriberName: s.name,
		amqpURL:        s.amqpURL,
		pattern:        pattern,
		errors:         make(chan error),
	}

	ch, rcv, err := sub.init()
	if err != nil {
		return nil, err
	}

	go func() {
		err := sub.receiveLoop(ch, rcv)
		if err != nil {
			sub.retryLoop()
		}
	}()

	return sub, nil
}
