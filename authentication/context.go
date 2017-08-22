package authentication

import (
	"context"
	"errors"
)

type ctxKey struct{}
type serviceTokenCtxKey struct{}

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

// GetServiceTokenFromCtx gets a previously stored service authentication token from the context
func GetServiceTokenFromCtx(ctx context.Context) string {
	t, ok := ctx.Value(serviceTokenCtxKey{}).(string)
	if !ok {
		panic("no service token in context")
	}

	return t
}

// SetServiceTokenInCtx sets a service authentication token in a new context and returns it
func SetServiceTokenInCtx(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, serviceTokenCtxKey{}, token)
}
