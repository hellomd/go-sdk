package recovery

import (
	"fmt"
	"net/http"

	raven "github.com/getsentry/raven-go"
	"github.com/hellomd/go-sdk/recovery/sentry"
)

// NewMiddleware creates a new middleware that recovers from errors with optional sentry integration
func NewMiddleware(sentryDSN string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		ctx := r.Context()

		cli, err := raven.New(sentryDSN)
		if err == nil {
			cli.SetHttpContext(raven.NewHttp(r))
			ctx = sentry.SetInCtx(r.Context(), cli)
		}

		next(w, r.WithContext(ctx))

		defer func() {
			if err := recover(); err != nil {
				HandleError(ctx, fmt.Errorf("%v", err))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
	}
}
