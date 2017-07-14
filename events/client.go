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

// NewClient creates a new client for publishing and consuming events
func NewClient(amqpURL string) (Client, error) {
	c := &client{amqpURL: amqpURL}
	if err := c.bootstrap(); err != nil {
		return nil, fmt.Errorf("error bootstrapping events: %v", err)
	}

	return c, nil
}

type client struct {
	amqpURL      string
	rlock, wlock sync.Mutex
	connection   *amqp.Connection
}

func (c *client) Publish(key string, body interface{}) error {
	ch, err := c.channel()
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

func (c *client) Subscribe(pattern string) (Subscription, error) {
	return newSubscription(c, pattern)
}

func (c *client) bootstrap() error {
	ch, err := c.channel()
	if err != nil {
		return fmt.Errorf("error opening AMQP channel: %v", err)
	}

	if err := ch.ExchangeDeclare(ExchangeName, amqp.ExchangeTopic, durable, autoDelete, internal, noWait, nil); err != nil {
		return fmt.Errorf("error declaring exchange: %v", err)
	}

	return nil
}

func (c *client) channel() (*amqp.Channel, error) {
	conn, err := rabbit.GetConnection(c.amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return ch, nil
}
