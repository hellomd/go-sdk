package requestid

import (
	"context"
	"errors"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

// ErrNoRequestIDInCtx -
var ErrNoRequestIDInCtx = errors.New("No request id in context")

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
func GetRequestIDFromContext(ctx context.Context) (string, error) {
	id, ok := ctx.Value(RequestIDcontextKey).(string)
	if !ok {
		return "", ErrNoRequestIDInCtx
	}
	return id, nil
}

func (mw *requestID) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var rID string
	if rID = r.Header.Get(reqIDheaderKey); rID == "" {
		rID = uuid.NewV4().String()
	}

	w.Header().Set(reqIDheaderKey, rID)
	next(w, r.WithContext(context.WithValue(r.Context(), RequestIDcontextKey, rID)))
}
