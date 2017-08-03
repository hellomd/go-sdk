package recovery

import (
	"context"
	"fmt"
	"net/http"

	raven "github.com/getsentry/raven-go"
	"github.com/hellomd/go-sdk/recovery/sentry"
)

// Middleware adds a sentry client to the context and recovers from errors using HandleError
type Middleware struct {
	SentryDSN   string
	HandleError func(context.Context, error)
}

// NewMiddleware returns a new recovery middleware with the default HandleError function
func NewMiddleware(SentryDSN string) *Middleware {
	return &Middleware{SentryDSN, HandleError}
}

// ServeHTTP Adds sentry client to context and recover from errors
func (mw *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := r.Context()

	cli, err := raven.New(mw.SentryDSN)
	if err == nil {
		cli.SetHttpContext(raven.NewHttp(r))
		ctx = sentry.SetInCtx(r.Context(), cli)
	}

	defer func() {
		if err := recover(); err != nil {
			mw.HandleError(ctx, fmt.Errorf("%v", err))
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	next(w, r.WithContext(ctx))
}
