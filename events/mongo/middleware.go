package mongo

import (
	"context"

	"github.com/hellomd/go-sdk/events"
	"github.com/hellomd/go-sdk/mongo"
)

// NewMiddleware -
func NewMiddleware(mongoURL string, useSSL bool) (func(context.Context, *events.Event, events.HandlerFunc), error) {
	mw, err := mongo.NewMiddleware(mongoURL, useSSL)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context, event *events.Event, next events.HandlerFunc) {
		mw.UseSession(ctx, func(ctx context.Context) {
			next(ctx, event)
		})
	}, nil
}
