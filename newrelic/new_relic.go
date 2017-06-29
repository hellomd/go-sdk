package newrelic

import (
	"net/http"

	newrelic "github.com/newrelic/go-agent"
)

// NewRelic -
type NewRelic interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type newRelic struct {
	newRelicApp newrelic.Application
}

// NewNewRelic -
func NewNewRelic(newRelicApp newrelic.Application) NewRelic {
	return &newRelic{newRelicApp}
}

func (mw *newRelic) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	txn := mw.newRelicApp.StartTransaction(r.URL.Path, w, r)
	defer txn.End()
	next(w, r)
}
