package events

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hellomd/go-sdk/rabbit"
	"github.com/streadway/amqp"
)

const (
	// ExchangeName is the exchange to which events are published
	ExchangeName = "x-events"

	durable    = true
	autoDelete = false
	internal   = false
	noWait     = false
	mandatory  = false
	immediate  = false
)

// NewPublisher creates a new client that can publish events
func NewPublisher(amqpURL string) (Publisher, error) {
	c := &publisher{amqpURL: amqpURL}
	if err := c.bootstrap(); err != nil {
		return nil, fmt.Errorf("error bootstrapping events: %v", err)
	}

	return c, nil
}

type publisher struct {
	amqpURL      string
	rlock, wlock sync.Mutex
	connection   *amqp.Connection
}

func (c *publisher) Publish(key string, body interface{}) error {
	ch, err := newChannel(c.amqpURL)
	if err != nil {
		return fmt.Errorf("error opening AMQP channel: %v", err)
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling body: %v", err)
	}

	pub := amqp.Publishing{Body: bodyJSON}
	if err := ch.Publish(ExchangeName, key, mandatory, immediate, pub); err != nil {
		return fmt.Errorf("error publishing event: %v", err)
	}

	return nil
}

func (c *publisher) bootstrap() error {
	ch, err := newChannel(c.amqpURL)
	if err != nil {
		return fmt.Errorf("error opening AMQP channel: %v", err)
	}

	if err := ch.ExchangeDeclare(ExchangeName, amqp.ExchangeTopic, durable, autoDelete, internal, noWait, nil); err != nil {
		return fmt.Errorf("error declaring exchange: %v", err)
	}

	return nil
}

func newChannel(amqpURL string) (*amqp.Channel, error) {
	conn, err := rabbit.GetConnection(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return ch, nil
}
