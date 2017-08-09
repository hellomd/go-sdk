package errors

import "net/http"

// NotFoundHandler -
type notFoundHandler struct{}

// NewNotFoundHandler -
func NewNotFoundHandler() http.Handler {
	return &notFoundHandler{}
}

// ServeHTTP -
func (nh *notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("route not found"))
}
