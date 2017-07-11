package etag

import (
	"crypto/md5"
	"fmt"
	"net/http"
)

const (
	//ETagHeaderKey -
	ETagHeaderKey = "ETag"

	//IfNonMatchHeaderKey -
	IfNonMatchHeaderKey = "If-None-Match"
)

type etagResponseWriter struct {
	http.ResponseWriter
	req  *http.Request
	code int
}

func (erw *etagResponseWriter) Write(b []byte) (int, error) {
	etag := etag(b)
	erw.Header().Set(ETagHeaderKey, etag)
	if erw.req.Header.Get(IfNonMatchHeaderKey) == etag {
		erw.WriteHeader(http.StatusNotModified)
		return erw.Write(nil)
	}

	erw.WriteHeader(erw.code)
	return erw.ResponseWriter.Write(b)
}

func etag(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}

func (erw *etagResponseWriter) WriteHeader(code int) {
	erw.code = code
	erw.ResponseWriter.WriteHeader(code)
}

// Middleware -
type Middleware interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type middleware struct{}

// NewMiddleware -
func NewMiddleware() Middleware {
	return &middleware{}
}

func (mw *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(&etagResponseWriter{w, r, http.StatusOK}, r)
}
