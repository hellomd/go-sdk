package logmatic

import (
	"time"

	logmatic "github.com/logmatic/logmatic-go"
	"github.com/sirupsen/logrus"
)

// Config is the configuration for the hook
type Config struct {
	// APIKey is your private key for using Logmatic API
	APIKey string

	// Interval is how often bufferred logs will be sent to Logmatic (default is 5 seconds)
	Interval time.Duration

	// BufferSize is the maximum number of events that will be sent to Logmatic at once (default is 100)
	BufferSize int
}

// NewLogrusHook creates a new hook for sending logs to Logmatic with Logrus
func NewLogrusHook(config Config) logrus.Hook {
	if config.APIKey == "" {
		panic("logmatic API key is required")
	}

	if config.Interval == 0 {
		config.Interval = 5 * time.Second
	}

	if config.BufferSize <= 0 {
		config.BufferSize = 100
	}

	return &hook{
		newAggregator(config.APIKey, config.BufferSize, config.Interval),
		&logmatic.JSONFormatter{},
	}
}

type hook struct {
	*aggregator
	formatter logrus.Formatter
}

func (h hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h hook) Fire(entry *logrus.Entry) error {
	json, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}

	h.aggregator.Write(string(json))
	return nil
}
