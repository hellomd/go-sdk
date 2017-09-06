package events

import (
	"context"
)

// Acknowledger sends feedback on whether the event has been successfully processed or not
type Acknowledger interface {
	Reject(bool)
	Ack()
}

// Event is what is received by a subscription
type Event struct {
	Acknowledger
	Key    string
	Body   []byte
	Header map[string]string
}

type Handler interface {
	Process(context.Context, *Event)
}

type HandlerFunc func(context.Context, *Event)
