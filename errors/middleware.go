package errors

import (
	"io/ioutil"
	"net/http"

	"bytes"

	"encoding/json"

	"github.com/hellomd/go-sdk/logger"
	"github.com/sirupsen/logrus"
)

// errorReponseWriter - wrapper to ResponseWriter
type errorReponseWriter struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func (erw *errorReponseWriter) Write(data []byte) (int, error) {
	if erw.status == 0 || erw.status >= 400 {
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
	//logger *logrus.Logger
}

// NewMiddleware -
func NewMiddleware() Middleware {
	return &middleware{}
}

func (mw *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	eWriter := &errorReponseWriter{w, http.StatusOK, bytes.NewBuffer([]byte{})}
	next(eWriter, r)

	entry, err := logger.GetFromCtx(r.Context())
	if err != nil {
		logger := logrus.New()
		logger.Out = ioutil.Discard
		entry = logrus.NewEntry(logger)
	}

	if eWriter.status >= http.StatusBadRequest && eWriter.status < http.StatusInternalServerError {
		entry.Warn(eWriter.body.String())
		if eWriter.status == http.StatusNotFound {
			ErrNotFound.Description = eWriter.body.String()
			json.NewEncoder(eWriter.ResponseWriter).Encode(ErrNotFound)
			return
		}
		eWriter.ResponseWriter.Write(eWriter.body.Bytes())
	}

	if eWriter.status >= http.StatusInternalServerError {
		entry.Error(eWriter.body.String())
		json.NewEncoder(eWriter.ResponseWriter).Encode(ErrUnexptectedError)
	}

}
