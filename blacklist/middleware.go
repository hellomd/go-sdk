package blacklist

import (
	"net/http"
)

// Middleware -
type Middleware interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type middleware struct {
}

// NewMiddleware -
func NewMiddleware() Middleware {
	return &middleware{}
}

// ServeHTTP -
func (mw *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.URL.Path == healthcheck {
		w.WriteHeader(http.StatusOK)
		return
	}

	for _, x := range blacklist {
		if r.URL.Path == x {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	next(w, r)
}
