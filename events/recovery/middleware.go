package recovery

import (
	"context"
	"fmt"
	"log"

	raven "github.com/getsentry/raven-go"
	"github.com/hellomd/go-sdk/events"
	"github.com/hellomd/go-sdk/logger"
	"github.com/hellomd/go-sdk/recovery/sentry"
)

type Recovery struct {
	sentryDSN   string
	handleError func(ctx context.Context, err error)
}

func NewMiddleware(sentryDSN string) *Recovery {
	return &Recovery{sentryDSN, handleError}
}

func (r *Recovery) Process(ctx context.Context, event *events.Event, next events.HandlerFunc) {
	cli, err := raven.New(r.sentryDSN)
	if err == nil {
		ctx = sentry.SetInCtx(ctx, cli)
	}

	defer func() {
		log.Println("lele")
		if rMsg := recover(); rMsg != nil {
			log.Println("lala")
			err := fmt.Errorf("panic recovered: %v", rMsg)
			event.Reject(false, err)
			r.handleError(ctx, err)
		}
	}()

	next(ctx, event)
}

// HandleError reports given error to logger and sentry when they are available
func handleError(ctx context.Context, err error) {
	logger, ctxErr := logger.GetFromCtx(ctx)
	if ctxErr == nil {
		logger.Error(err)
	}

	sentry, ctxErr := sentry.GetFromCtx(ctx)
	if ctxErr == nil {
		sentry.CaptureError(err, nil)
	}
}
