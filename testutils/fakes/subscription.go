package fakes

import (
	"sync"

	"github.com/hellomd/go-sdk/events"
)

// Subscription should be used to simulate receiving a subscribed event
type Subscription struct {
	lock    sync.Mutex
	receive chan events.Event
	close   chan struct{}
}

// Receive returns a channel to which fake subscription events are passed
func (s *Subscription) Receive() <-chan events.Event {
	s.assert()
	return s.receive
}

// NotifyClose notifies of when the subscription ends
func (s *Subscription) NotifyClose() <-chan struct{} {
	s.assert()
	return s.close
}

// Close stops receiving events for this subscription
func (s *Subscription) Close() {
	s.assert()
	close(s.receive)
	s.close <- struct{}{}
}

// Errors does nothing in this fake implementation
func (s *Subscription) Errors() <-chan error {
	s.assert()
	return make(chan error)
}

// FakeReceive sends an event to the Receive channel
func (s *Subscription) FakeReceive(key string, body []byte) {
	s.assert()
	s.receive <- events.Event{
		Key:  key,
		Body: body,
	}
}

func (s *Subscription) assert() {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.receive == nil {
		s.receive = make(chan events.Event)
	}

	if s.close == nil {
		s.close = make(chan struct{})
	}
}

var _ subscription = &Subscription{}
var _ subscription = &events.Subscription{}

type subscription interface {
	Receive() <-chan events.Event
	NotifyClose() <-chan struct{}
	Close()
	Errors() <-chan error
}
