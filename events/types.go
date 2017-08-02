package events

// Acknowledger sends feedback on whether the event has been successfully processed or not
type Acknowledger interface {
	Reject(bool)
	Ack()
}

// Event is what is received by a subscription
type Event struct {
	Acknowledger
	Key  string
	Body []byte
}
