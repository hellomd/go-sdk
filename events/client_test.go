package events

import (
	"sync"
	"testing"
	"time"

	"github.com/hellomd/go-sdk/config"
)

func TestPublishSubscribe(t *testing.T) {
	client, err := NewClient(config.Get("AMQP_URL"))
	if err != nil {
		t.Error(err)
		return
	}

	sub, err := client.Subscribe("questions.*.created")
	if err != nil {
		t.Error(err)
		return
	}

	defer sub.Close()

	messages := [][]byte{}
	wg := new(sync.WaitGroup)
	go func() {
		wg.Add(1)
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

	if err := client.Publish("questions.article.created", map[string]string{"foo": "bar"}); err != nil {
		t.Error(err)
	}

	if err := client.Publish("questions.product.created", map[string]int{"one": 1}); err != nil {
		t.Error(err)
	}

	wg.Wait()

	if len(messages) != 2 {
		t.Errorf("expected to receive 2 messages, but got %v", len(messages))
		return
	}

	if string(messages[0]) != `{"foo":"bar"}` {
		t.Errorf("expected first message to be %v, but was %v", `{"foo":"bar"}`, string(messages[0]))
	}

	if string(messages[1]) != `{"one":1}` {
		t.Errorf("expected first message to be %v, but was %v", `{"one":1}`, string(messages[1]))
	}
}
