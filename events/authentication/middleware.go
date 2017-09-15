package authentication

import (
	"context"

	"github.com/hellomd/go-sdk/authentication"
	"github.com/hellomd/go-sdk/events"
)

func NewMiddleware(secret []byte) func(context.Context, *events.Event, events.HandlerFunc) {
	ctxAuth := authentication.NewContextMiddleware(secret)
	return func(ctx context.Context, event *events.Event, next events.HandlerFunc) {
		ctx, err := ctxAuth(ctx, event.Header[authentication.HeaderKey])
		if err != nil {
			event.Reject(false, err)
			return
		}

		next(ctx, event)
	}
}
