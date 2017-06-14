package middlewares

import (
	"net/http"

	"fmt"

	raven "github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
)

// Recovery -
type Recovery interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

// RavenClient -
type RavenClient interface {
	CaptureError(error, map[string]string, ...raven.Interface) string
	SetHttpContext(*raven.Http)
}

type recovery struct {
	RavenClient
	*logrus.Logger
}

// NewRecovery -
func NewRecovery(ravenClient RavenClient, logger *logrus.Logger) Recovery {
	return &recovery{ravenClient, logger}
}

func (mw *recovery) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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
