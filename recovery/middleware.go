package recovery

import (
	"net/http"

	"fmt"

	raven "github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
)

// Middleware -
type Middleware interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

// RavenClient -
type RavenClient interface {
	CaptureError(error, map[string]string, ...raven.Interface) string
	SetHttpContext(*raven.Http)
}

type middleware struct {
	RavenClient
	*logrus.Logger
}

// NewMiddleware -
func NewMiddleware(ravenClient RavenClient, logger *logrus.Logger) Middleware {
	return &middleware{ravenClient, logger}
}

func (mw *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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
