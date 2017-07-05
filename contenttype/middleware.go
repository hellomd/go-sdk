package contenttype

import (
	"net/http"
)

const (
	contentTypeHeaderKey = "Content-Type"
	defaultContentType   = "application/json"
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
	w.Header().Set(contentTypeHeaderKey, defaultContentType)
	next(w, r)
}
