package newrelic

import (
	"net/http"

	newrelic "github.com/newrelic/go-agent"
)

// Middleware -
type Middleware interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type middleware struct {
	newRelicApp newrelic.Application
}

// NewMiddleware -
func NewMiddleware(newRelicApp newrelic.Application) Middleware {
	return &middleware{newRelicApp}
}

func (mw *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	txn := mw.newRelicApp.StartTransaction(r.URL.Path, w, r)
	defer txn.End()
	next(w, r)
}
