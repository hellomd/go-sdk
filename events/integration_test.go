package events

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/hellomd/go-sdk/config"
	"github.com/hellomd/go-sdk/rabbit"
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
			}
		}
	}()

	if err := publisher.Publish("questions.article.created", map[string]string{"foo": "bar"}); err != nil {
		t.Error(err)
	}

	if err := publisher.Publish("questions.product.created", map[string]int{"one": 1}); err != nil {
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

func TestReconnection(t *testing.T) {
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

	// publish first event
	if err := publisher.Publish("questions.article.created", "one"); err != nil {
		t.Error(err)
	}

	var evt Event

	// get first event
	select {
	case evt = <-sub.Receive():
		if string(evt.Body) != `"one"` {
			t.Errorf(`expected "one", but got %v`, string(evt.Body))
			return
		}
	case <-time.After(50 * time.Millisecond):
		t.Errorf("expected to receive a message, but got none")
		return
	}

	// force disconnection
	{
		conn, err := rabbit.GetConnection(amqpURL)
		if err != nil {
			t.Errorf("unexpected connection error: %v", err)
			return
		}

		if err := conn.Close(); err != nil {
			t.Errorf("failed to force connection to close: %v", err)
			return
		}
	}

	time.Sleep(10 * time.Millisecond)

	// publish new event
	if err := publisher.Publish("questions.product.created", "two"); err != nil {
		t.Error(err)
	}

	// get new event
	select {
	case evt = <-sub.Receive():
		if string(evt.Body) != `"two"` {
			t.Errorf(`expected "two", but got %v`, string(evt.Body))
			return
		}
	case <-time.After(50 * time.Millisecond):
		t.Errorf("expected to receive a message, but got none")
		return
	}
}
