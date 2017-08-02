package recovery

import (
	"context"
)

type ctxKey struct{}

// GetReporterFromCtx -
func GetReporterFromCtx(ctx context.Context) (*Reporter, error) {
	reporter, ok := ctx.Value(ctxKey{}).(*Reporter)
	if !ok {
		return nil, errNotInCtx
	}
	return reporter, nil
}

// SetReporterInCtx -
func SetReporterInCtx(ctx context.Context, reporter *Reporter) context.Context {
	return context.WithValue(ctx, ctxKey{}, reporter)
}
