package middlewares

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const realIPHeaderKey = "X-Real-IP"

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

// Logger -
type Logger interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type logger struct {
	logger *logrus.Logger
}

// NewLogger -
func NewLogger(l *logrus.Logger) Logger {
	return &logger{l}
}

func (mw *logger) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	lw := newLoggerReponseWriter(w)
	next(lw, r)

	entry := logrus.NewEntry(mw.logger)

	remoteAddr := r.RemoteAddr
	if realIP := r.Header.Get(realIPHeaderKey); realIP != "" {
		remoteAddr = realIP
	}

	latency := time.Since(start)
	requestID := r.Context().Value(RequestIDcontextKey)

	entry.WithFields(logrus.Fields{
		"request_id": requestID,
		"path":       r.RequestURI,
		"method":     r.Method,
		"remote":     remoteAddr,
		"took":       latency,
		"status":     lw.status,
	}).Info("")
}
