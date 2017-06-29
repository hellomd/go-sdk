package requestid

import (
	"context"
	"errors"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

type requestIDCtxKey int

const (
	// RequestIDCtxKey -
	RequestIDCtxKey requestIDCtxKey = iota

	// RequestIDHeaderKey -
	RequestIDHeaderKey = "X-Request-ID"
)

// ErrNoRequestIDInCtx -
var ErrNoRequestIDInCtx = errors.New("No request id in context")

// Middleware -
type Middleware interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type middleware struct{}

// NewMiddleware -
func NewMiddleware() Middleware {
	return &middleware{}
}

// GetRequestIDFromCtx -
func GetRequestIDFromCtx(ctx context.Context) (string, error) {
	id, ok := ctx.Value(RequestIDCtxKey).(string)
	if !ok {
		return "", ErrNoRequestIDInCtx
	}
	return id, nil
}

func (mw *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var rID string
	if rID = r.Header.Get(RequestIDHeaderKey); rID == "" {
		rID = uuid.NewV4().String()
	}

	w.Header().Set(RequestIDHeaderKey, rID)
	next(w, r.WithContext(context.WithValue(r.Context(), RequestIDCtxKey, rID)))
}
