package etag

import (
	"crypto/md5"
	"fmt"
	"net/http"
)

const (
	//ETagHeaderKey -
	ETagHeaderKey = "ETag"

	//IfNoneMatchHeaderKey -
	IfNoneMatchHeaderKey = "If-None-Match"
)

type etagResponseWriter struct {
	http.ResponseWriter
	req         *http.Request
	code        int
	wroteHeader bool
}

func (erw *etagResponseWriter) Write(b []byte) (int, error) {
	etag := etag(b)
	erw.Header().Set(ETagHeaderKey, etag)
	if erw.req.Header.Get(IfNoneMatchHeaderKey) == etag {
		erw.ResponseWriter.WriteHeader(http.StatusNotModified)
		return erw.ResponseWriter.Write(nil)
	}
	if erw.wroteHeader == false {
		erw.wroteHeader = true
		erw.ResponseWriter.WriteHeader(erw.code)
	}
	return erw.ResponseWriter.Write(b)
}

func etag(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}

func (erw *etagResponseWriter) WriteHeader(code int) {
	erw.code = code
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
	writer := &etagResponseWriter{w, r, http.StatusOK, false}
	next(writer, r)
	if !writer.wroteHeader {
		writer.ResponseWriter.WriteHeader(writer.code)
	}
}
