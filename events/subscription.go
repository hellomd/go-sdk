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

func newSubscription(client *client, pattern string) (*subscription, error) {
	s := &subscription{
		client:  client,
		pattern: pattern,
		errors:  make(chan error),
	}

	ch, rcv, err := s.init()
	if err != nil {
		return nil, err
	}

	go func() {
		err := s.receiveLoop(ch, rcv)
		if err != nil {
			s.retryLoop()
		}
	}()

	return s, nil
}

type subscription struct {
	client  *client
	pattern string
	errors  chan error

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

	randomID := "87fsadh"
	queueName := "q-sub-" + s.pattern + randomID

	const (
		durable    = false
		autoDelete = true
		exclusive  = true
		noWait     = false
		autoAck    = true
		noLocal    = false
	)

	ch, err := s.client.channel()
	if err != nil {
		return nil, nil, fmt.Errorf("error opening AMQP channel: %v", err)
	}

	if _, err := ch.QueueDeclare(queueName, durable, autoDelete, exclusive, noWait, nil); err != nil {
		return nil, nil, fmt.Errorf("error declaring queue: %v", err)
	}

	if err := ch.QueueBind(queueName, s.pattern, ExchangeName, noWait, nil); err != nil {
		return nil, nil, fmt.Errorf("error binding subscription queue: %v", err)
	}

	rcv, err := ch.Consume(queueName, randomID, autoAck, exclusive, noLocal, noWait, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating consumer: %v", err)
	}

	return ch, rcv, nil
}

func (s *subscription) receiveLoop(ch *amqp.Channel, rcv <-chan amqp.Delivery) error {
	chCloser := ch.NotifyClose(make(chan *amqp.Error))

	for {
		select {
		case delivery := <-rcv: // send event
			s.events <- delivery.Body
		case <-s.closer: // stop receiving
			ch.Close()
			return nil
		case err := <-chCloser: // reconnect
			return err
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
