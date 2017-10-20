package events

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hellomd/go-sdk/rabbit"
	"github.com/streadway/amqp"
)

// ExchangeName is the exchange to which events are published
const ExchangeName = "x-events"

const (
	durable    = true
	autoDelete = false
	internal   = false
	noWait     = false
	mandatory  = false
	immediate  = false
)

// NewPublisher creates a new client that can publish events
//
// This publisher is concurrent safe and should be reused as much as possible
// because of its initialization logic that involves declaring the RabbitMQ exchange.
func NewPublisher(amqpURL string) (*Publisher, error) {
	return NewPublisherCustom(amqpURL, ExchangeName)
}

// NewPublisherCustom creates a new client that can publish events
//
// This publisher is concurrent safe and should be reused as much as possible
// because of its initialization logic that involves declaring the RabbitMQ exchange.
func NewPublisherCustom(amqpURL, exchange string) (*Publisher, error) {
	c := &Publisher{amqpURL: amqpURL, Exchange: exchange}
	if err := c.bootstrap(); err != nil {
		return nil, fmt.Errorf("error bootstrapping events: %v", err)
	}

	return c, nil
}

// Publisher is a client that can publish events
type Publisher struct {
	Exchange     string
	amqpURL      string
	rlock, wlock sync.Mutex
}

// Publish publishes an event with default headers
func (c *Publisher) Publish(ctx context.Context, key string, body interface{}) error {
	h, err := DefaultHeaders(ctx)
	if err != nil {
		return err
	}

	return c.PublishH(key, body, h)
}

// PublishH publishes an event with custom headers
func (c *Publisher) PublishH(key string, body interface{}, headers map[string]string) error {
	ch, err := newChannel(c.amqpURL)
	if err != nil {
		return fmt.Errorf("error opening AMQP channel: %v", err)
	}

	defer ch.Close()

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling body: %v", err)
	}

	headersTable := amqp.Table{}
	for k, v := range headers {
		headersTable[k] = v
	}

	pub := amqp.Publishing{
		Body:    bodyJSON,
		Headers: headersTable,
	}

	if err := ch.Publish(c.Exchange, key, mandatory, immediate, pub); err != nil {
		return fmt.Errorf("error publishing event: %v", err)
	}

	return nil
}

func (c *Publisher) bootstrap() error {
	ch, err := newChannel(c.amqpURL)
	if err != nil {
		return fmt.Errorf("error opening AMQP channel: %v", err)
	}

	defer ch.Close()

	if err := ch.ExchangeDeclare(c.Exchange, amqp.ExchangeTopic, durable, autoDelete, internal, noWait, nil); err != nil {
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
