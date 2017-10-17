package requestid

import (
	"context"

	"github.com/hellomd/go-sdk/events"
	"github.com/hellomd/go-sdk/requestid"
)

func NewMiddleware() events.PipelineHandlerFunc {
	return func(ctx context.Context, event *events.Event, next events.HandlerFunc) {
		if id, ok := event.Header[requestid.RequestIDHeaderKey]; ok {
			ctx = context.WithValue(ctx, requestid.RequestIDCtxKey, id)
		}

		next(ctx, event)
	}
}
