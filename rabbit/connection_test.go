package rabbit

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hellomd/go-sdk/config"
	"github.com/streadway/amqp"
)

var amqpURL = config.Get("AMQP_URL")

func TestMain(m *testing.M) {
	timeout := time.Now().Add(5 * time.Second)
	for {
		conn, err := amqp.Dial(amqpURL)
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
