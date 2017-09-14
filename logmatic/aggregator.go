package logmatic

import (
	"sync"
	"time"
)

type aggregator struct {
	client     *client
	entries    chan string
	interval   time.Duration
	lock       sync.Mutex
	hasStarted bool
}

func newAggregator(apiKey string, bufferSize int, interval time.Duration) *aggregator {
	a := &aggregator{
		client:     newClient(apiKey),
		entries:    make(chan string, bufferSize),
		interval:   interval,
		lock:       sync.Mutex{},
		hasStarted: false,
	}
	a.start()
	return a
}

func (a *aggregator) Write(event string) {
	for {
		select {
		case a.entries <- event:
			return
		default:
			a.flush()
		}
	}
}

func (a *aggregator) start() {
	if a.hasStarted {
		panic("logmatic aggregator has already been started")
	}
	a.hasStarted = true

	go func() {
		for {
			<-time.After(a.interval)
			a.flush()
		}
	}()
}

func (a *aggregator) flush() {
	a.lock.Lock()
	defer a.lock.Unlock()

	flushed := flushChan(a.entries)
	go a.client.Send(flushed)
}

func flushChan(events <-chan string) []string {
	maxFlushed := len(events)
	flushed := make([]string, 0, maxFlushed)
	for i := 0; i < maxFlushed; i++ {
		select {
		case ev := <-events:
			flushed = append(flushed, ev)
		default:
			return flushed
		}
	}
	return flushed
}
