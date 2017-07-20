package authentication

import (
	"context"
	"errors"
)

type ctxKey int

// CtxKey -
const CtxKey ctxKey = 0

// ErrNoCurrentUserInCtx -
var ErrNoCurrentUserInCtx = errors.New("no current user in context")

// GetCurrentUserFromCtx -
func GetCurrentUserFromCtx(ctx context.Context) (*CurrentUser, error) {
	db, ok := ctx.Value(CtxKey).(*CurrentUser)
	if !ok {
		return nil, ErrNoCurrentUserInCtx
	}
	return db, nil
}

// SetCurrentUserOnCtx -
func SetCurrentUserOnCtx(ctx context.Context, user *CurrentUser) context.Context {
	return context.WithValue(ctx, CtxKey, user)
}
