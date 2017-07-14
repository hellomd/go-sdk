package etag

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
)

const (
	//ETagHeaderKey -
	ETagHeaderKey = "ETag"

	//IfNoneMatchHeaderKey -
	IfNoneMatchHeaderKey = "If-None-Match"
)

var (
	errDoubleWrite = errors.New("Write called twice with ETag Middleware")
)

type etagResponseWriter struct {
	http.ResponseWriter
	req         *http.Request
	code        int
	wroteHeader bool
	wroteBody   bool
}

func (erw *etagResponseWriter) Write(b []byte) (int, error) {
	etag := calculateEtag(b)
	erw.Header().Set(ETagHeaderKey, etag)

	if erw.req.Header.Get(IfNoneMatchHeaderKey) == etag {
		erw.code = http.StatusNotModified
		erw.writeHeader()
		return erw.writeBody(nil)
	}

	erw.writeHeader()
	return erw.writeBody(b)
}

func calculateEtag(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}

func (erw *etagResponseWriter) writeBody(data []byte) (int, error) {
	if erw.wroteBody {
		panic(errDoubleWrite)
	}
	erw.wroteBody = true
	return erw.ResponseWriter.Write(data)
}

func (erw *etagResponseWriter) writeHeader() {
	if !erw.wroteHeader {
		erw.wroteHeader = true
		erw.ResponseWriter.WriteHeader(erw.code)
	}
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
	if r.Method == http.MethodGet {
		writer := &etagResponseWriter{w, r, http.StatusOK, false, false}
		next(writer, r)
		if !writer.wroteHeader {
			writer.ResponseWriter.WriteHeader(writer.code)
		}
	} else {
		next(w, r)
	}
}
