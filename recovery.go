package middlewares

import (
	"net/http"

	"fmt"

	raven "github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
)

// RecoveryMiddleware -
type RecoveryMiddleware interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

// RavenClient -
type RavenClient interface {
	CaptureError(error, map[string]string, ...raven.Interface) string
	SetHttpContext(*raven.Http)
}

type recoveryMiddleware struct {
	RavenClient
	*logrus.Logger
}

// NewRecoveryMiddleware -
func NewRecoveryMiddleware(ravenClient RavenClient, logger *logrus.Logger) RecoveryMiddleware {
	return &recoveryMiddleware{ravenClient, logger}
}

func (mw *recoveryMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			mw.SetHttpContext(raven.NewHttp(r))
			mw.CaptureError(fmt.Errorf("%v", err), nil)
			mw.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	next(w, r)
}
