package rabbit

import (
	"testing"

	"sync"

	"time"

	"github.com/hellomd/go-sdk/config"
)

var amqpURL = config.Get("AMQP_URL")

func TestGetNewConnection(t *testing.T) {
	conn, err := GetConnection(amqpURL)
	if err != nil {
		t.Error(err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		t.Error(err)
	}

	ch.Close()
}

func TestGetSameConnection(t *testing.T) {
	conn, err := GetConnection(amqpURL)
	if err != nil {
		t.Error(err)
	}

	defer conn.Close()

	conn2, err := GetConnection(amqpURL)
	if err != nil {
		t.Error(err)
	}

	if conn != conn2 {
		t.Error("expected conn and conn2 to be the same")
	}
}

func TestConnectionConcurrency(t *testing.T) {
	// arbitrary milliseconds to wait for openning and closing each connection
	presets := map[int]int{
		0:   500,
		25:  300,
		12:  600,
		700: 100,
	}

	wg := &sync.WaitGroup{}
	// do the following 4000 times (1000 for each of the 4 delay presets)
	for i := 0; i < 1000; i++ {
		for start, end := range presets {
			wg.Add(1)
			go func(start, end time.Duration) {
				defer wg.Done()

				time.Sleep(start)

				conn, err := GetConnection(amqpURL)
				if err != nil {
					t.Error(err)
					return
				}

				defer conn.Close()

				time.Sleep(end)
			}(time.Duration(start)*time.Millisecond, time.Duration(end)*time.Millisecond)
		}
	}
	wg.Wait()
}
