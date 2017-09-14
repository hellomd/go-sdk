# Logmatic

A Logrus hook that sends logs to Logmatic via HTTP.

The hook sends events in batch, based on an `Interval` and a `BufferSize`. The events will only be sent
at every tick of `Interval`, or when the number of events reaches `BufferSize`.

## Usage

```go
package main

import (
	"time"

	"github.com/hellomd/go-sdk/logmatic"
	"github.com/sirupsen/logrus"
)

func main() {
	hook := logmatic.NewLogrusHook(logmatic.Config{
		APIKey:     "your-api-key",
		Interval:   5 * time.Second, // how often bufferred logs will be sent to Logmatic (default is 5 seconds)
		BufferSize: 100,             // maximum number of events that will be sent to Logmatic at once (default is 100)
	})

	log := logrus.New()
	log.Hooks.Add(hook)
}
```
