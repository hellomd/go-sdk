package events

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/hellomd/go-sdk/config"
	"github.com/streadway/amqp"
)

func TestMain(m *testing.M) {
	timeout := time.Now().Add(5 * time.Second)
	for {
		conn, err := amqp.Dial(config.Get("AMQP_URL"))
		if err != nil && time.Now().After(timeout) {
			fmt.Println("Failed to get AMQP connection:", err)
			os.Exit(1)
		} else if err != nil {
			continue
		}

		conn.Close()
		break
	}
	os.Exit(m.Run())
}

func TestPublishSubscribe(t *testing.T) {
	amqpURL := config.Get("AMQP_URL")

	publisher, err := NewPublisher(amqpURL)
	if err != nil {
		t.Error(err)
		return
	}

	subscriber := NewSubscriber("testsub", amqpURL)

	sub, err := subscriber.Subscribe("questions.*.created")
	if err != nil {
		t.Error(err)
		return
	}

	defer sub.Close()

	messages := []Event{}
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		timeout := time.After(1 * time.Second)

		for {
			select {
			case <-timeout:
				return

			case evt := <-sub.Receive():
				messages = append(messages, evt)
				evt.Ack()
			}
		}
	}()

	if err := publisher.Publish("questions.article.created", map[string]string{"foo": "bar"}, nil); err != nil {
		t.Error(err)
	}

	if err := publisher.Publish("questions.product.created", map[string]int{"one": 1}, nil); err != nil {
		t.Error(err)
	}

	wg.Wait()

	if len(messages) != 2 {
		t.Errorf("expected to receive 2 messages, but got %v", len(messages))
		return
	}

	if string(messages[0].Body) != `{"foo":"bar"}` {
		t.Errorf("expected first message to be %v, but was %v", `{"foo":"bar"}`, string(messages[0].Body))
	}

	if string(messages[1].Body) != `{"one":1}` {
		t.Errorf("expected first message to be %v, but was %v", `{"one":1}`, string(messages[1].Body))
	}
}
