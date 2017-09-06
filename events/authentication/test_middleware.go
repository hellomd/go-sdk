package authentication

import (
	"context"

	"github.com/hellomd/go-sdk/authentication"
	"github.com/hellomd/go-sdk/events"
)

type TestMiddleware struct {
	user *authentication.User
}

func NewTestMiddleware() *TestMiddleware {
	return &TestMiddleware{&authentication.User{}}
}

func (t *TestMiddleware) SetUser(user *authentication.User) {
	t.user = user
}

func (t *TestMiddleware) Process(ctx context.Context, event *events.Event, next events.HandlerFunc) {
	ctx = authentication.SetUserInCtx(ctx, t.user)
	ctx = authentication.SetServiceTokenInCtx(ctx, "fake-service-token")
	next(ctx, event)
}
