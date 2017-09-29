package blacklist

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/urfave/negroni"
)

func TestValidRequest(t *testing.T) {
	called := false
	mw := NewMiddleware().(*middleware)

	req := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	a := negroni.New()
	a.Use(mw)
	a.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		called = true
		next(w, r)
	}))

	a.ServeHTTP(response, req)
	if !called {
		t.Errorf("Expected middleware to call next")
	}
}

func TestHealthcheckRequest(t *testing.T) {
	notCalled := true
	mw := NewMiddleware().(*middleware)

	req := httptest.NewRequest("GET", "/healthcheck", nil)
	response := httptest.NewRecorder()

	a := negroni.New()
	a.Use(mw)
	a.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		notCalled = false
		next(w, r)
	}))

	a.ServeHTTP(response, req)
	if !notCalled {
		t.Errorf("Expected middleware not to be called next")
	}

	if response.Code != http.StatusOK {
		t.Errorf("Status should be OK for /healthcheck")
	}
}

func TestBlacklistedRequest(t *testing.T) {
	mw := NewMiddleware().(*middleware)

	for _, x := range blacklist {
		notCalled := true
		req := httptest.NewRequest("GET", x, nil)
		response := httptest.NewRecorder()

		a := negroni.New()
		a.Use(mw)
		a.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			notCalled = false
			next(w, r)
		}))

		a.ServeHTTP(response, req)
		if !notCalled {
			t.Errorf("Expected middleware not to be called next")
		}

		if response.Code != http.StatusNotFound {
			t.Errorf("Status should be Not Found for %v request", x)
		}
	}

}
