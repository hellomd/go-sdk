package authentication

import (
	"context"
	"errors"
)

type ctxKey struct{}

var errNoUserInCtx = errors.New("no user in context")

// GetUserFromCtx -
func GetUserFromCtx(ctx context.Context) *User {
	u, ok := ctx.Value(ctxKey{}).(*User)
	if !ok {
		panic(errNoUserInCtx)
	}
	return u
}

// SetUserInCtx -
func SetUserInCtx(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, ctxKey{}, user)
}
