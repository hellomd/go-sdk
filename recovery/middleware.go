package recovery

import (
	"fmt"
	"net/http"

	raven "github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
)

// NewMiddleware returns a middleware that recovers from errors and adds a error reporter to context
func NewMiddleware(sentryDSN string, logger *logrus.Logger) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		cli, err := raven.New(sentryDSN)
		if err != nil {
			cli = nil
		}

		reporter := &Reporter{cli, logger, raven.NewHttp(r)}
		ctx := SetReporterInCtx(r.Context(), reporter)
		defer func() {
			if err := recover(); err != nil {
				reporter.Error(fmt.Errorf("%v", err))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
	}
}
