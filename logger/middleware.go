package logger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hellomd/go-sdk/errors"
	"github.com/hellomd/go-sdk/requestid"
	"github.com/sirupsen/logrus"
)

// RealIPHeaderKey -
const RealIPHeaderKey = "X-Real-IP"

// loggerReponseWriter - wrapper to ResponseWriter
type loggerReponseWriter struct {
	http.Flusher
	http.ResponseWriter
	http.CloseNotifier
	status int
	body   string
}

func newLoggerReponseWriter(w http.ResponseWriter) *loggerReponseWriter {
	var flusher http.Flusher
	var cNotifier http.CloseNotifier
	var ok bool
	if flusher, ok = w.(http.Flusher); !ok {
		flusher = nil
	}

	if cNotifier, ok = w.(http.CloseNotifier); !ok {
		cNotifier = nil
	}

	return &loggerReponseWriter{flusher, w, cNotifier, http.StatusOK, ""}
}

func (lrw *loggerReponseWriter) Write(body []byte) (int, error) {
	if lrw.status >= http.StatusBadRequest {
		lrw.body = string(body)
	}
	return lrw.ResponseWriter.Write(body)
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
	appName     string
	environment string
	logger      *logrus.Logger
}

// NewMiddleware -
func NewMiddleware(appName, environment string, logger *logrus.Logger) Middleware {
	return &middleware{appName, environment, logger}
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
		"request_id":       requestID,
		"remote":           remoteAddr,
		"path":             r.RequestURI,
		"method":           r.Method,
		"application_name": mw.appName,
		"environment":      mw.environment,
	})

	lw := newLoggerReponseWriter(w)
	next(lw, r.WithContext(SetInCtx(r.Context(), entry)))

	latency := time.Since(start)
	message := fmt.Sprintf("%v %v | %v | %v \"%v\"", r.Method, r.RequestURI, latency, lw.status, lw.GetErrorMessage())
	newEntry := entry.WithFields(logrus.Fields{
		"took":   latency / time.Millisecond,
		"status": lw.status,
	})

	if lw.status >= http.StatusBadRequest && lw.status < http.StatusInternalServerError {
		newEntry.Warn(message)
		return
	}

	if lw.status >= http.StatusInternalServerError {
		newEntry.Error(message)
		return
	}

	newEntry.Info(message)
}

// GetErrorMessage -
func (lrw *loggerReponseWriter) GetErrorMessage() string {
	if lrw.status >= http.StatusBadRequest {
		jsonError := errors.JSONError{}
		err := json.Unmarshal([]byte(lrw.body), &jsonError)
		if err != nil {
			return lrw.body
		}
		return jsonError.Message
	}
	return http.StatusText(lrw.status)
}
