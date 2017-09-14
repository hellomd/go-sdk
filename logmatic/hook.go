package logmatic

import (
	"time"

	logmatic "github.com/logmatic/logmatic-go"
	"github.com/sirupsen/logrus"
)

type Config struct {
	APIKey     string
	Interval   time.Duration
	BufferSize int
}

func NewLogrusHook(config Config) logrus.Hook {
	if config.APIKey == "" {
		panic("logmatic API key should is required")
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
