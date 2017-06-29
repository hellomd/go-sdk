package requestid

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/urfave/negroni"
)

func TestWithoutRequestId(t *testing.T) {
	response := httptest.NewRecorder()
	a := negroni.New()
	a.Use(NewMiddleware())
	a.ServeHTTP(response, httptest.NewRequest("GET", "/", nil))
	if response.HeaderMap.Get(RequestIDHeaderKey) == "" {
		t.Errorf("Expected some value in header %s, but is empty", RequestIDHeaderKey)
	} else {
		t.Log(response.HeaderMap.Get(RequestIDHeaderKey))
	}
}

func TestWithRequestId(t *testing.T) {
	response := httptest.NewRecorder()
	a := negroni.New()
	a.Use(NewMiddleware())

	req := httptest.NewRequest("GET", "/", nil)
	myReqID := "42"
	req.Header.Set(RequestIDHeaderKey, myReqID)

	a.ServeHTTP(response, req)

	if v := response.HeaderMap.Get(RequestIDHeaderKey); v != myReqID {
		t.Errorf("Expected '%s', but got '%s' in header %s.", myReqID, v, RequestIDHeaderKey)
	} else {
		t.Log(v)
	}
}

func TestContextSet(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	myReqID := "42"
	req.Header.Set(RequestIDHeaderKey, myReqID)

	response := httptest.NewRecorder()
	a := negroni.New()
	a.Use(NewMiddleware())
	a.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if v := r.Context().Value(RequestIDCtxKey); v != myReqID {
			t.Errorf("Expected '%s', but got '%s' in header %s.", myReqID, v, RequestIDHeaderKey)
		} else {
			t.Log(v)
		}
		next(w, r)
	}))

	a.ServeHTTP(response, req)

}

func TestGetRequestIDFromCtx(t *testing.T) {
	ctx := context.Background()
	_, err := GetRequestIDFromCtx(ctx)
	if err != ErrNoRequestIDInCtx {
		t.Error("Expected ErrNoRequestIDInCtx, got: ", err)
	}

	expectedReqID := uuid.NewV4().String()
	ctx = context.WithValue(ctx, RequestIDCtxKey, expectedReqID)
	requestID, _ := GetRequestIDFromCtx(ctx)
	if requestID != expectedReqID {
		t.Errorf("Expected %s, got: %s", expectedReqID, requestID)
	}
}
