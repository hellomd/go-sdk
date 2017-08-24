package events_test

import (
	"encoding/json"
	"fmt"

	"github.com/hellomd/go-sdk/events"
	"github.com/sirupsen/logrus"
)

func ExampleSubscriber() {
	subscriber := events.NewSubscriber("feed", "amqp://guest:guest@localhost.com", logrus.New())

	sub, err := subscriber.Subscribe("question.*.created")
	if err != nil {
		panic(fmt.Errorf("error subscribing to questions created: %v", err))
	}

	// Log errors
	go func() {
		for {
			select {
			case err := <-sub.Errors():
				fmt.Println("Event error:", err)

			case <-sub.NotifyClose():
				return
			}
		}
	}()

	type question struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}

	// Consume events
	for {
		select {
		case evt := <-sub.Receive():
			q := question{}
			if err := json.Unmarshal(evt.Body, &q); err != nil {
				fmt.Println("Error deserializing event:", err)
				continue
			}

			fmt.Println("Question", q.ID, "was created")

		case <-sub.NotifyClose():
			return
		}
	}
}
