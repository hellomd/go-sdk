package recovery

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/hellomd/go-sdk/logger"
	"github.com/hellomd/go-sdk/recovery/sentry"
)

// HandleError reports given error to logger and sentry when they are available
func HandleError(ctx context.Context, err error) {
	logger, ctxErr := logger.GetFromCtx(ctx)
	if ctxErr == nil {
		logger.Error(err)
	}

	sentry, ctxErr := sentry.GetFromCtx(ctx)
	if ctxErr == nil {
		sentry.CaptureError(err, nil)
	}
}

// DebugHandleError reports given error to logger including stacktrace and sentry when they are available
func DebugHandleError(ctx context.Context, err error) {
	logger, ctxErr := logger.GetFromCtx(ctx)
	if ctxErr == nil {
		logger.WithField("stack_trace", fmt.Sprintf("%s", debug.Stack())).Error(err)
	}

	sentry, ctxErr := sentry.GetFromCtx(ctx)
	if ctxErr == nil {
		sentry.CaptureError(err, nil)
	}
}
