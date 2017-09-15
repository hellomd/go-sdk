package logger

import (
	"context"
	"fmt"
	"time"

	"github.com/hellomd/go-sdk/events"
	"github.com/hellomd/go-sdk/logger"
	"github.com/hellomd/go-sdk/requestid"
	"github.com/sirupsen/logrus"
)

const (
	statusAcked    = "acked"
	statusRequeued = "requeued"
	statusRejected = "rejected"
)

var _ events.Acknowledger = &loggerAcknowledger{}

func NewMiddleware(appName, env string, instance *logrus.Logger) events.PipelineHandlerFunc {
	return func(ctx context.Context, e *events.Event, next events.HandlerFunc) {
		start := time.Now()
		requestID := ctx.Value(requestid.RequestIDCtxKey)

		entry := logrus.NewEntry(instance)
		entry = entry.WithFields(logrus.Fields{
			"request_id":       requestID,
			"event_key":        e.Key,
			"application_name": appName,
			"environment":      env,
		})

		ack := &loggerAcknowledger{inner: e.Acknowledger}
		ew := *e
		ew.Acknowledger = ack

		next(logger.SetInCtx(ctx, entry), &ew)

		latency := time.Since(start)
		message := fmt.Sprintf("evt %v | %v | %v", e.Key, latency, ack.status)
		entry = entry.WithFields(logrus.Fields{
			"took":   latency,
			"status": ack.status,
		})

		switch ack.status {
		case statusRequeued:
			entry.WithField("error", fmt.Sprint(ack.err)).Warn(fmt.Sprintf(`%v "%v"`, message, ack.err))

		case statusRejected:
			entry.WithField("error", fmt.Sprint(ack.err)).Error(fmt.Sprintf(`%v "%v"`, message, ack.err))

		default:
			entry.Info(message)
		}
	}
}

type loggerAcknowledger struct {
	inner  events.Acknowledger
	status string
	err    error
}

func (a *loggerAcknowledger) Reject(requeue bool, err error) {
	a.inner.Reject(requeue, err)
	if a.status != "" {
		return
	}

	a.err = err

	if requeue {
		a.status = statusRequeued
		return
	}

	a.status = statusRejected
}

func (a *loggerAcknowledger) Ack() {
	a.inner.Ack()
	if a.status != "" {
		return
	}

	a.status = statusAcked
}
