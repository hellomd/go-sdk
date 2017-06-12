package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/urfave/negroni"
)

func TestWithoutRequestId(t *testing.T) {
	response := httptest.NewRecorder()
	a := negroni.New()
	a.Use(NewRequestID())
	a.ServeHTTP(response, httptest.NewRequest("GET", "/", nil))
	if response.HeaderMap.Get(headerKey) == "" {
		t.Errorf("Expected some value in header %s, but is empty", headerKey)
	} else {
		t.Log(response.HeaderMap.Get(headerKey))
	}
}

func TestWithRequestId(t *testing.T) {
	response := httptest.NewRecorder()
	a := negroni.New()
	a.Use(NewRequestID())

	req := httptest.NewRequest("GET", "/", nil)
	myReqID := "42"
	req.Header.Set(headerKey, myReqID)

	a.ServeHTTP(response, req)

	if v := response.HeaderMap.Get(headerKey); v != myReqID {
		t.Errorf("Expected '%s', but got '%s' in header %s.", myReqID, v, headerKey)
	} else {
		t.Log(v)
	}
}

func TestContextSet(t *testing.T) {
	//Prepare Request
	req := httptest.NewRequest("GET", "/", nil)
	myReqID := "42"
	req.Header.Set(headerKey, myReqID)

	response := httptest.NewRecorder()
	a := negroni.New()
	a.Use(NewRequestID())
	a.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if v := r.Context().Value(RequestIDcontextKey); v != myReqID {
			t.Errorf("Expected '%s', but got '%s' in header %s.", myReqID, v, headerKey)
		} else {
			t.Log(v)
		}
		next(w, r)
	}))

	a.ServeHTTP(response, req)

}
