package authentication

import (
	"context"
	"errors"
)

type ctxKey int

// CtxKey -
const CtxKey ctxKey = 0

// errUserInCtx -
var errNoUserInCtx = errors.New("no user in context")

// GetUserFromCtx -
func GetUserFromCtx(ctx context.Context) *User {
	u, ok := ctx.Value(CtxKey).(*User)
	if !ok {
		panic(errNoUserInCtx)
	}
	return u
}

// SetUserInCtx -
func SetUserInCtx(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, CtxKey, user)
}
