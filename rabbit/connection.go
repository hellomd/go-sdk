package rabbit

import (
	"sync"

	"github.com/streadway/amqp"
)

var connections = map[string]*amqp.Connection{}
var rlock, wlock sync.Mutex

// GetConnection gets an open AMQP connection for the given URL, creating it if necessary
func GetConnection(url string) (*amqp.Connection, error) {
	rlock.Lock()
	defer rlock.Unlock()

	if c, ok := connections[url]; ok {
		return c, nil
	}

	return dial(url)
}

func dial(url string) (*amqp.Connection, error) {
	wlock.Lock()
	defer wlock.Unlock()

	if conn, ok := connections[url]; ok {
		return conn, nil
	}

	c, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	go func() {
		<-c.NotifyClose(make(chan *amqp.Error))
		rlock.Lock()
		defer rlock.Unlock()

		wlock.Lock()
		defer wlock.Unlock()

		delete(connections, url)
	}()

	connections[url] = c
	return c, nil
}
