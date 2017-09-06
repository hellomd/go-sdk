package events_test

import (
	"fmt"

	"github.com/hellomd/go-sdk/events"
)

func ExamplePublisher() {
	pub, err := events.NewPublisher("amqp://guest:guest@localhost")
	if err != nil {
		panic(fmt.Errorf("there was a problem asserting AMQP structure: %v", err))
	}

	err = pub.Publish("questions.article.created", map[string]interface{}{
		"id":    "abc123",
		"title": "What is the meaning of life?",
	}, nil)
	if err != nil {
		panic(fmt.Errorf("error publishing question: %v", err))
	}
}
