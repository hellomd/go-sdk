package logger

import (
	"net/http"
	"time"

	"github.com/hellomd/middlewares/requestid"
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

	lw := newLoggerReponseWriter(w)
	next(lw, r)

	entry := logrus.NewEntry(mw.logger)

	remoteAddr := r.RemoteAddr
	if realIP := r.Header.Get(RealIPHeaderKey); realIP != "" {
		remoteAddr = realIP
	}

	latency := time.Since(start)
	requestID := r.Context().Value(requestid.RequestIDCtxKey)

	entry.WithFields(logrus.Fields{
		"request_id": requestID,
		"path":       r.RequestURI,
		"method":     r.Method,
		"remote":     remoteAddr,
		"took":       latency,
		"status":     lw.status,
	}).Info("")
}
