package events

import "context"
import "github.com/sirupsen/logrus"

const (
	baseRetryDelay     float64 = 10
	maxRetryDelay      float64 = 5000
	defaultConcurrency         = 5
)

// NewSubscriber creates a new client that can subscribe to events
//
// This subscriber is concurrent safe and has no initialization logic.
// Feel free to use whatever lifecycle you think is best for it.
func NewSubscriber(name, amqpURL string, logger *logrus.Logger) *Subscriber {
	return &Subscriber{name, amqpURL, logger, defaultConcurrency, ExchangeName}
}

// Subscriber is a client that can subscribe to events
type Subscriber struct {
	name        string
	amqpURL     string
	logger      *logrus.Logger
	Concurrency int
	Exchange    string
}

// Subscribe creates a new subscription for receiving events with keys that match the given pattern
func (s *Subscriber) Subscribe(pattern string) (*Subscription, error) {
	sub := &Subscription{
		exchange:       s.Exchange,
		subscriberName: s.name,
		amqpURL:        s.amqpURL,
		pattern:        pattern,
		errors:         make(chan error),
		concurrency:    int(float64(s.Concurrency) * 1.2),
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

// SubscribeH creates a new subscription for receiving events with keys that match the given pattern
func (s *Subscriber) SubscribeH(pattern string, handler Handler) error {
	sub, err := s.Subscribe(pattern)
	if err != nil {
		return err
	}

	go func() {
		closer := sub.NotifyClose()
		receiver := sub.Receive()
		for {
			select {
			case <-closer:
				return

			case event := <-receiver:
				ctx := context.Background()
				go func() {
					defer func() {
						if r := recover(); r != nil {
							s.logger.WithFields(logrus.Fields{
								"error":        r,
								"subscription": pattern,
							}).Error("Subscription handler panicked")
						}
					}()
					handler.Process(ctx, &event)
				}()
			}
		}
	}()

	return nil
}
