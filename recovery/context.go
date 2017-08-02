package recovery

import (
	"context"
)

type ctxKey struct{}

// GetReporterFromCtx get the reporter from given context, returns error if not present
func GetReporterFromCtx(ctx context.Context) (*Reporter, error) {
	reporter, ok := ctx.Value(ctxKey{}).(*Reporter)
	if !ok {
		return nil, errNotInCtx
	}
	return reporter, nil
}

// SetReporterInCtx sets given reporter to given context
func SetReporterInCtx(ctx context.Context, reporter *Reporter) context.Context {
	return context.WithValue(ctx, ctxKey{}, reporter)
}