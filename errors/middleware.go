package errors

import (
	"net/http"

	"bytes"

	"github.com/sirupsen/logrus"
)

// errorReponseWriter - wrapper to ResponseWriter
type errorReponseWriter struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func (erw *errorReponseWriter) Write(data []byte) (int, error) {
	if erw.status == 0 || erw.status >= 500 {
		return erw.body.Write(data)
	}
	return erw.ResponseWriter.Write(data)
}

func (erw *errorReponseWriter) WriteHeader(code int) {
	erw.status = code
	erw.ResponseWriter.WriteHeader(code)
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
	eWriter := &errorReponseWriter{w, http.StatusOK, bytes.NewBuffer([]byte{})}
	next(eWriter, r)

	entry := logrus.NewEntry(mw.logger)

	if eWriter.status == http.StatusInternalServerError {
		entry.Error(eWriter.body.String())
		eWriter.ResponseWriter.Write([]byte(http.StatusText(http.StatusInternalServerError)))
	}

}
