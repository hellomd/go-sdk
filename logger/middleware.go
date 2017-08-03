package logger

import (
	"net/http"
	"time"

	"github.com/hellomd/go-sdk/requestid"
	"github.com/sirupsen/logrus"
)

// RealIPHeaderKey -
const RealIPHeaderKey = "X-Real-IP"

// loggerReponseWriter - wrapper to ResponseWriter
type loggerReponseWriter struct {
	http.ResponseWriter
	status int
}

func newLoggerReponseWriter(w http.ResponseWriter) *loggerReponseWriter {
	return &loggerReponseWriter{w, http.StatusOK}
}

func (lrw *loggerReponseWriter) WriteHeader(code int) {
	lrw.status = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Middleware -
type Middleware interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type middleware struct {
	logger *logrus.Logger
}

// NewMiddleware -
func NewMiddleware(l *logrus.Logger) Middleware {
	return &middleware{l}
}

func (mw *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	requestID := r.Context().Value(requestid.RequestIDCtxKey)

	remoteAddr := r.RemoteAddr
	if realIP := r.Header.Get(RealIPHeaderKey); realIP != "" {
		remoteAddr = realIP
	}

	entry := logrus.NewEntry(mw.logger)
	entry = entry.WithFields(logrus.Fields{
		"request_id": requestID,
		"remote":     remoteAddr,
	})

	lw := newLoggerReponseWriter(w)
	next(lw, r.WithContext(SetInCtx(r.Context(), entry)))

	latency := time.Since(start)

	entry.WithFields(logrus.Fields{
		"path":   r.RequestURI,
		"method": r.Method,
		"took":   latency,
		"status": lw.status,
	}).Info("")
}
