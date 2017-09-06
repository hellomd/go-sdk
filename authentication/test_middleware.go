package authentication

import (
	"net/http"
)

// TestMiddleware -
type TestMiddleware struct {
	user *User
}

// NewTestMiddleware -
func NewTestMiddleware() *TestMiddleware {
	return &TestMiddleware{&User{}}
}

// SetUser -
func (t *TestMiddleware) SetUser(user *User) {
	t.user = user
}

func (t *TestMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := r.Context()
	ctx = SetUserInCtx(ctx, t.user)
	ctx = SetServiceTokenInCtx(ctx, "fake-service-token")
	next(w, r.WithContext(ctx))
}
