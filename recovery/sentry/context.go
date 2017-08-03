package sentry

import (
	"context"
	"errors"

	raven "github.com/getsentry/raven-go"
)

var errNotInCtx = errors.New("no sentry client in context")

type ctxKey struct{}

// GetFromCtx gets a sentry client from context
func GetFromCtx(ctx context.Context) (*raven.Client, error) {
	sentry, ok := ctx.Value(ctxKey{}).(*raven.Client)
	if !ok {
		return nil, errNotInCtx
	}
	return sentry, nil
}

// SetInCtx sets a given sentry client to given context
func SetInCtx(ctx context.Context, sentry *raven.Client) context.Context {
	return context.WithValue(ctx, ctxKey{}, sentry)
}
