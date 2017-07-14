package etag

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
	"sync"
)

const (
	//ETagHeaderKey -
	ETagHeaderKey = "ETag"

	//IfNoneMatchHeaderKey -
	IfNoneMatchHeaderKey = "If-None-Match"
)

var (
	errDoubleWrite = errors.New("write called twice with ETag Middleware")
)

type etagResponseWriter struct {
	http.ResponseWriter
	req         *http.Request
	code        int
	wroteHeader bool
	wroteBody   bool
	headerLock  sync.Mutex
	bodyLock    sync.Mutex
}

func (erw *etagResponseWriter) Write(b []byte) (int, error) {
	etag := calculateEtag(b)
	erw.Header().Set(ETagHeaderKey, etag)

	if erw.req.Header.Get(IfNoneMatchHeaderKey) == etag {
		if erw.req.Method == http.MethodGet {
			erw.code = http.StatusNotModified
		}
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
	erw.bodyLock.Lock()
	defer erw.bodyLock.Unlock()

	if erw.wroteBody {
		panic(errDoubleWrite)
	}
	erw.wroteBody = true
	return erw.ResponseWriter.Write(data)
}

func (erw *etagResponseWriter) writeHeader() {
	erw.headerLock.Lock()
	defer erw.headerLock.Unlock()

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
	writer := &etagResponseWriter{w, r, http.StatusOK, false, false, sync.Mutex{}, sync.Mutex{}}
	next(writer, r)
	writer.writeHeader()
}
