package contenttype

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/urfave/negroni"
)

func TestMiddleware(t *testing.T) {
	mw := NewMiddleware().(*middleware)
	req := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	a := negroni.New()
	a.Use(mw)
	a.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.Write([]byte("{}"))
		next(w, r)
	}))

	a.ServeHTTP(response, req)
	contentType := response.Header().Get(contentTypeHeaderKey)
	if contentType != defaultContentType {
		t.Errorf("Expected Content-Type %s, got %s", defaultContentType, contentType)
	}
}

func TestMiddlewareRespectsPreviouslySet(t *testing.T) {
	expectedContentType := "text/plain"
	mw := NewMiddleware().(*middleware)
	req := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	a := negroni.New()
	a.Use(mw)
	a.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.Header().Set(contentTypeHeaderKey, expectedContentType)
		next(w, r)
	}))

	a.ServeHTTP(response, req)
	contentType := response.Header().Get(contentTypeHeaderKey)
	if contentType != expectedContentType {
		t.Errorf("Expected Content-Type %s, got %s", expectedContentType, contentType)
	}
}

func TestMiddlewareDoesNothingOnEmptyBody(t *testing.T) {
	mw := NewMiddleware().(*middleware)
	req := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	a := negroni.New()
	a.Use(mw)
	a.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		next(w, r)
	}))

	a.ServeHTTP(response, req)
	contentType := response.Header().Get(contentTypeHeaderKey)
	if contentType != "" {
		t.Errorf("Expected no Content-Type got %s", contentType)
	}
}
