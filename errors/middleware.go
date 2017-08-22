package errors

import (
	"net/http"
	"strings"

	"bytes"

	"encoding/json"
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
}

// NewMiddleware -
func NewMiddleware() Middleware {
	return &middleware{}
}

func (mw *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	eWriter := &errorReponseWriter{w, http.StatusOK, bytes.NewBuffer([]byte{})}
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
