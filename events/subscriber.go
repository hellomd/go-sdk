package events

import (
	"fmt"
	"math"

	"time"

	"github.com/streadway/amqp"
)

const (
	baseRetryDelay float64 = 10
	maxRetryDelay  float64 = 5000
)

// NewSubscriber creates a new client that can subscribe to events
//
// This subscriber is concurrent safe and has no initialization logic.
// Feel free to use whatever lifecycle you think is best for it.
func NewSubscriber(name, amqpURL string) Subscriber {
	return &subscriber{name, amqpURL}
}

type subscriber struct {
	name    string
	amqpURL string
}

func (s *subscriber) Subscribe(pattern string) (Subscription, error) {
	sub := &subscription{
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

type subscription struct {
	subscriberName string
	amqpURL        string
	pattern        string
	errors         chan error

	events     chan []byte
	closer     chan struct{}
	retryCount float64
}

func (s *subscription) Receive() <-chan []byte {
	return s.events
}

func (s *subscription) NotifyClose() <-chan struct{} {
	return s.closer
}

func (s *subscription) Close() {
	close(s.closer)
}

func (s *subscription) Errors() <-chan error {
	return s.errors
}

func (s *subscription) init() (*amqp.Channel, <-chan amqp.Delivery, error) {
	s.events = make(chan []byte)
	s.closer = make(chan struct{})

	queueName := fmt.Sprintf("q-sub-%v-%v", s.subscriberName, s.pattern)

	const (
		durable    = true
		autoDelete = false
		exclusive  = false
		noWait     = false
		autoAck    = true
		noLocal    = false
	)

	ch, err := newChannel(s.amqpURL)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening AMQP channel: %v", err)
	}

	if _, err := ch.QueueDeclare(queueName, durable, autoDelete, exclusive, noWait, nil); err != nil {
		return nil, nil, fmt.Errorf("error declaring queue: %v", err)
	}

	if err := ch.QueueBind(queueName, s.pattern, ExchangeName, noWait, nil); err != nil {
		return nil, nil, fmt.Errorf("error binding subscription queue: %v", err)
	}

	rcv, err := ch.Consume(queueName, s.subscriberName, autoAck, exclusive, noLocal, noWait, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating consumer: %v", err)
	}

	return ch, rcv, nil
}

func (s *subscription) receiveLoop(ch *amqp.Channel, rcv <-chan amqp.Delivery) error {
	chCloser := ch.NotifyClose(make(chan *amqp.Error))

	for {
		select {
		case err := <-chCloser: // reconnect
			return err
		case delivery := <-rcv: // send event
			if delivery.Acknowledger != nil { // safeguard against closed channel sends
				s.events <- delivery.Body
			}
		case <-s.closer: // stop receiving
			ch.Close()
			return nil
		}
	}
}

func (s *subscription) retryLoop() {
	for {
		err := s.retry()
		if err != nil {
			s.retryCount++
			delay := time.Duration(
				math.Min(
					baseRetryDelay*(s.retryCount*s.retryCount),
					maxRetryDelay,
				)) * time.Millisecond
			s.logError(fmt.Errorf("retrying in %v: %v", delay, err))
			time.Sleep(delay)
			continue
		}

		return
	}
}

func (s *subscription) retry() error {
	ch, rcv, err := s.init()
	if err != nil {
		return err
	}

	s.retryCount = 0
	if err := s.receiveLoop(ch, rcv); err != nil {
		return err
	}

	return nil
}

func (s *subscription) logError(err error) {
	select {
	case s.errors <- err:
	default:
	}
}
