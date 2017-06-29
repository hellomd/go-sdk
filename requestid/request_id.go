package requestid

import (
	"context"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

type contextKey int

const (
	// RequestIDcontextKey -
	RequestIDcontextKey contextKey = iota
	reqIDheaderKey                 = "X-Request-ID"
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

// GetRequestIDFromContext -
func GetRequestIDFromContext(ctx context.Context) string {
	id, ok := ctx.Value(RequestIDcontextKey).(string)
	if !ok {
		panic("Could not lookup requestID from context")
	}
	return id
}

func (mw *requestID) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var rID string
	if rID = r.Header.Get(reqIDheaderKey); rID == "" {
		rID = uuid.NewV4().String()
	}

	w.Header().Set(reqIDheaderKey, rID)
	next(w, r.WithContext(context.WithValue(r.Context(), RequestIDcontextKey, rID)))
}
