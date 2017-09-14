package logmatic

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	defaultEndpoint = "https://api.logmatic.io/v1/input/"
	maxAttempts     = 5
	retryInterval   = 500 * time.Millisecond
)

type client struct {
	apiKey     string
	httpClient *http.Client
}

func newClient(apiKey string) *client {
	return &client{apiKey, &http.Client{Timeout: 5 * time.Second}}
}

func (c *client) Send(events []string) {
	i := 1
	for {
		if err := c.TrySend(events); err != nil && i < maxAttempts {
			log.Println("error sending logs, will retry in", retryInterval, ":", err)
			time.Sleep(retryInterval)
			i++
			continue
		} else if err != nil {
			log.Println("error sending logs, will no longer retry:", err)
		}

		break
	}
}

func (c *client) TrySend(events []string) error {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Println("Logmatic panic:", rec)
		}
	}()

	if len(events) == 0 {
		return nil
	}

	body := bytes.NewBufferString("[" + strings.Join(events, ",") + "]")
	req, _ := http.NewRequest("POST", defaultEndpoint+c.apiKey, body)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %v sending events", res.StatusCode)
	}

	return nil
}
