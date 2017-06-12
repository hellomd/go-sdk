package middlewares

import (
	"context"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

type contextKey int

const (
	// RequestIDcontextKey -
	RequestIDcontextKey contextKey = iota
	headerKey                      = "X-Request-ID"
)

// RequestID -
type RequestID interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type requestID struct{}

// NewRequestID -
func NewRequestID() RequestID {
	return &requestID{}
}

func (mw *requestID) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var rID string
	if rID = r.Header.Get(headerKey); rID == "" {
		rID = uuid.NewV4().String()
	}

	w.Header().Set(headerKey, rID)
	next(w, r.WithContext(context.WithValue(r.Context(), RequestIDcontextKey, rID)))
}
