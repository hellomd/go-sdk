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
	if response.HeaderMap.Get(reqIDheaderKey) == "" {
		t.Errorf("Expected some value in header %s, but is empty", reqIDheaderKey)
	} else {
		t.Log(response.HeaderMap.Get(reqIDheaderKey))
	}
}

func TestWithRequestId(t *testing.T) {
	response := httptest.NewRecorder()
	a := negroni.New()
	a.Use(NewRequestID())

	req := httptest.NewRequest("GET", "/", nil)
	myReqID := "42"
	req.Header.Set(reqIDheaderKey, myReqID)

	a.ServeHTTP(response, req)

	if v := response.HeaderMap.Get(reqIDheaderKey); v != myReqID {
		t.Errorf("Expected '%s', but got '%s' in header %s.", myReqID, v, reqIDheaderKey)
	} else {
		t.Log(v)
	}
}

func TestContextSet(t *testing.T) {
	//Prepare Request
	req := httptest.NewRequest("GET", "/", nil)
	myReqID := "42"
	req.Header.Set(reqIDheaderKey, myReqID)

	response := httptest.NewRecorder()
	a := negroni.New()
	a.Use(NewRequestID())
	a.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if v := r.Context().Value(RequestIDcontextKey); v != myReqID {
			t.Errorf("Expected '%s', but got '%s' in header %s.", myReqID, v, reqIDheaderKey)
		} else {
			t.Log(v)
		}
		next(w, r)
	}))

	a.ServeHTTP(response, req)

}
