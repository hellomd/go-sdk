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

	exptectedEtag := calculateEtag([]byte(data))

	if response.Header().Get(ETagHeaderKey) != exptectedEtag {
		t.Errorf("Unexpected ETag. Want: %v, got: %v", exptectedEtag, response.Header().Get(ETagHeaderKey))
	}
}

func TestResponseNotModified(t *testing.T) {
	srv := negroni.New(NewMiddleware())
	response := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add(IfNoneMatchHeaderKey, calculateEtag([]byte(data)))

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

func TestResponseNotModifiedOnlyOnGet(t *testing.T) {
	for _, method := range []string{"POST", "DELETE", "PUT", "PATCH", "HEAD"} {
		srv := negroni.New(NewMiddleware())
		response := httptest.NewRecorder()

		req := httptest.NewRequest(method, "/", nil)

		srv.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			w.Write([]byte(data))
			next(w, r)
		})

		srv.ServeHTTP(response, req)

		if response.Code != http.StatusOK {
			t.Errorf("Unexpected Status Code. Want: %v, got: %v", http.StatusOK, response.Code)
		}
	}
}

func TestPanicsOnSecondWrite(t *testing.T) {
	var thrownError error
	srv := negroni.New(NewMiddleware())
	response := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/", nil)

	srv.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		w.Write([]byte(data))
		defer func() {
			if r := recover(); r != nil {
				thrownError = r.(error)
			}
		}()
		w.Write([]byte(data))
		next(w, r)
	})

	srv.ServeHTTP(response, req)

	if thrownError != errDoubleWrite {
		t.Errorf("Expected double write to panic with error %v, got: %v", errDoubleWrite, thrownError)
	}
}

func TestDoesNotSetEtagWhenStatusCodeGt300(t *testing.T) {
	srv := negroni.New(NewMiddleware())
	response := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	srv.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(data))
		next(w, r)
	})

	srv.ServeHTTP(response, req)

	if etagHeader := response.Header().Get(ETagHeaderKey); etagHeader != "" {
		t.Errorf("Expected empty etag header, got: %v", etagHeader)
	}
}
