package events

import (
	"fmt"
	"math"
	"time"

	"github.com/streadway/amqp"
)

// Subscription can receive events to which it subscribes
type Subscription struct {
	exchange       string
	maxPriority    uint8
	subscriberName string
	amqpURL        string
	pattern        string
	errors         chan error

	events      chan Event
	closer      chan struct{}
	retryCount  float64
	concurrency int
}

// Receive returns a channel to which subscription events are passed
func (s *Subscription) Receive() <-chan Event {
	return s.events
}

// NotifyClose notifies of when the subscription is ended
func (s *Subscription) NotifyClose() <-chan struct{} {
	return s.closer
}

// Close stops receiving events for this subscription
func (s *Subscription) Close() {
	close(s.closer)
}

// Errors returns a channel where subscription errors are logged
func (s *Subscription) Errors() <-chan error {
	return s.errors
}

func (s *Subscription) init() (*amqp.Channel, <-chan amqp.Delivery, error) {
	s.events = make(chan Event)
	s.closer = make(chan struct{})

	queueName := fmt.Sprintf("q-sub-%v-%v", s.subscriberName, s.pattern)

	const (
		durable    = true
		autoDelete = false
		exclusive  = false
		noWait     = false
		autoAck    = false
		noLocal    = false
	)

	ch, err := newChannel(s.amqpURL)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening AMQP channel: %v", err)
	}

	err = ch.ExchangeDeclare(s.exchange, amqp.ExchangeTopic, durable, autoDelete, internal, noWait, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error declaring exchange: %s", err)
	}

	var queueArgs amqp.Table
	if s.maxPriority > 0 {
		queueArgs = amqp.Table{"x-max-priority": s.maxPriority}
	}
	if _, err := ch.QueueDeclare(queueName, durable, autoDelete, exclusive, noWait, queueArgs); err != nil {
		return nil, nil, fmt.Errorf("error declaring queue: %v", err)
	}

	if err := ch.QueueBind(queueName, s.pattern, s.exchange, noWait, nil); err != nil {
		return nil, nil, fmt.Errorf("error binding subscription queue: %v", err)
	}

	if err := ch.Qos(s.concurrency, 0, false); err != nil {
		return nil, nil, fmt.Errorf("error setting QoS: %v", err)
	}

	rcv, err := ch.Consume(queueName, s.subscriberName, autoAck, exclusive, noLocal, noWait, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating consumer: %v", err)
	}

	return ch, rcv, nil
}

func (s *Subscription) receiveLoop(ch *amqp.Channel, rcv <-chan amqp.Delivery) error {
	chCloser := ch.NotifyClose(make(chan *amqp.Error))

	for {
		select {
		case err := <-chCloser: // reconnect
			return err
		case delivery := <-rcv: // send event
			if delivery.Acknowledger != nil { // safeguard against closed channel sends
				header := map[string]string{}
				for k, v := range delivery.Headers {
					if sv, ok := v.(string); ok {
						header[k] = sv
					}
				}

				s.events <- Event{
					Acknowledger: newAcknowledger(&delivery),
					Key:          delivery.RoutingKey,
					Body:         delivery.Body,
					Header:       header,
				}
			}
		case <-s.closer: // stop receiving
			ch.Close()
			return nil
		}
	}
}

func (s *Subscription) retryLoop() {
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

func (s *Subscription) retry() error {
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

func (s *Subscription) logError(err error) {
	select {
	case s.errors <- err:
	default:
	}
}

func newAcknowledger(delivery *amqp.Delivery) *acknowledger {
	return &acknowledger{delivery}
}

type acknowledger struct {
	delivery *amqp.Delivery
}

func (a *acknowledger) Reject(requeue bool, _ error) {
	a.delivery.Reject(requeue)
}

func (a *acknowledger) Ack() {
	a.delivery.Ack(false)
}
