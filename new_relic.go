package middlewares

import (
	"net/http"

	newrelic "github.com/newrelic/go-agent"
)

// NewRelicMiddleware -
type NewRelicMiddleware interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type newRelicMiddleware struct {
	newRelicApp newrelic.Application
}

// NewNewRelic -
func NewNewRelic(newRelicApp newrelic.Application) NewRelicMiddleware {
	return &newRelicMiddleware{newRelicApp}
}

func (mw *newRelicMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	txn := mw.newRelicApp.StartTransaction(r.URL.Path, w, r)
	defer txn.End()
	next(w, r)
}
