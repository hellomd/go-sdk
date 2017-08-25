package errors

import (
	"net/http"
	"strings"

	"bytes"

	"encoding/json"
)

// errorReponseWriter - wrapper to ResponseWriter
type errorReponseWriter struct {
	http.Flusher
	http.ResponseWriter
	http.CloseNotifier
	status int
	body   *bytes.Buffer
}

func newErrorReponseWriter(w http.ResponseWriter) *errorReponseWriter {
	var flusher http.Flusher
	var cNotifier http.CloseNotifier
	var ok bool
	if flusher, ok = w.(http.Flusher); !ok {
		flusher = nil
	}

	if cNotifier, ok = w.(http.CloseNotifier); !ok {
		cNotifier = nil
	}

	return &errorReponseWriter{flusher, w, cNotifier, http.StatusOK, bytes.NewBuffer([]byte{})}
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
}

// NewMiddleware -
func NewMiddleware() Middleware {
	return &middleware{}
}

func (mw *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	eWriter := newErrorReponseWriter(w)
	next(eWriter, r)

	if eWriter.status >= http.StatusBadRequest && eWriter.status < http.StatusInternalServerError {
		if isErrorJSON(eWriter.body.Bytes()) {
			eWriter.ResponseWriter.Write(eWriter.body.Bytes())
		} else {
			resp := &JSONError{Code: errorCode(eWriter.status), Message: string(eWriter.body.Bytes())}
			eWriter.Header().Set("Content-Type", "application/json")
			json.NewEncoder(eWriter.ResponseWriter).Encode(resp)
		}
	}

	if eWriter.status >= http.StatusInternalServerError {
		json.NewEncoder(eWriter.ResponseWriter).Encode(&JSONError{Code: errorCode(eWriter.status), Message: "internal server error"})
	}
}

func isErrorJSON(body []byte) bool {
	return json.Unmarshal(body, &JSONError{}) == nil
}

func errorCode(status int) string {
	lowerCode := strings.ToLower(http.StatusText(status))
	return strings.Replace(lowerCode, " ", "_", -1)
}
