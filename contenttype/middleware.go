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

type contentTypeResponseWriter struct {
	http.ResponseWriter
}

func (erw *contentTypeResponseWriter) Write(b []byte) (int, error) {
	if b != nil && erw.ResponseWriter.Header().Get(contentTypeHeaderKey) == "" {
		erw.ResponseWriter.Header().Set(contentTypeHeaderKey, defaultContentType)
	} else {
		erw.ResponseWriter.Header().Set(contentTypeHeaderKey, http.DetectContentType(b))
	}
	return erw.ResponseWriter.Write(b)
}

// ServeHTTP -
func (mw *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(&contentTypeResponseWriter{w}, r)
}
