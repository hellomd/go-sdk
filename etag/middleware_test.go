package etag

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/urfave/negroni"
)

const (
	data = "MichaelDouglas"
)

func TestResponseWithEtag(t *testing.T) {
	srv := negroni.New(NewMiddleware())
	response := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/", nil)

	srv.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.Write([]byte(data))
		next(w, r)
	})

	srv.ServeHTTP(response, req)

	exptectedEtag := etag([]byte(data))

	if response.Header().Get(ETagHeaderKey) != exptectedEtag {
		t.Errorf("Unexpected ETag. Want: %v, got: %v", exptectedEtag, response.Header().Get(ETagHeaderKey))
	}
}

func TestResponseNotModified(t *testing.T) {
	srv := negroni.New(NewMiddleware())
	response := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add(IfNoneMatchHeaderKey, etag([]byte(data)))

	srv.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.Write([]byte(data))
		next(w, r)
	})

	srv.ServeHTTP(response, req)

	if response.Code != http.StatusNotModified {
		t.Errorf("Unexpected Status Code. Want: %v, got: %v", http.StatusNotModified, response.Code)
	}

	if len(response.Body.Bytes()) > 0 {
		t.Errorf("Unexpected Body. Want to be nil, got: %v", response.Body.Bytes())
	}
}
