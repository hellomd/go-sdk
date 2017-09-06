package fakes

import (
	"sync"

	"encoding/json"
	"fmt"

	"github.com/hellomd/go-sdk/events"
)

// Publisher should be used to fake an event publisher
type Publisher struct {
	lock      sync.Mutex
	published []events.Event
}

// Publish stores in the publisher instance the events that were published by it
func (p *Publisher) Publish(key string, body interface{}, header map[string]string) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.published == nil {
		p.published = []events.Event{}
	}

	marshalled, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling body: %v", err)
	}

	p.published = append(p.published, events.Event{Key: key, Body: marshalled, Header: header})
	return nil
}

// GetPublished returns all the events that were published to the fake publisher instance
func (p *Publisher) GetPublished() []events.Event {
	return p.published[:]
}

var _ publisher = &Publisher{}
var _ publisher = &events.Publisher{}

type publisher interface {
	Publish(key string, body interface{}, header map[string]string) error
}
