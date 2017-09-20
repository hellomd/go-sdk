package errors

import (
	"errors"
	"testing"

	"net/http"

	"net/http/httptest"

	"encoding/json"

	"reflect"

	raven "github.com/getsentry/raven-go"
	"github.com/urfave/negroni"
)

func TestError500(t *testing.T) {
	var actualErr error
	originalErr := errors.New("chaos")

	captureError := func(err error, tags map[string]string, interfaces ...raven.Interface) string {
		actualErr = err
		return ""
	}

	srv := negroni.New(&middleware{captureError})

	srv.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		http.Error(w, originalErr.Error(), http.StatusInternalServerError)
	})

	response := httptest.NewRecorder()

	srv.ServeHTTP(response, httptest.NewRequest("GET", "/", nil))

	resp := &JSONError{}
	err := json.NewDecoder(response.Body).Decode(resp)
	if err != nil {
		t.Errorf(`Unexpcted error on parse json, got :"%v" `, err.Error())
	}

	expectedError := &JSONError{Code: "internal_server_error", Message: "internal server error"}
	if !reflect.DeepEqual(resp, expectedError) {
		t.Errorf(`Unexpcted error response. Got :"%v", want: %v `, resp, expectedError)
	}

	if actualErr.Error() != originalErr.Error() {
		t.Errorf(`Unexpcted error captured. Got :"%v", want: %v `, actualErr, originalErr)
	}
}

func TestError422(t *testing.T) {

	srv := negroni.New(NewMiddleware())

	srv.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		info := BasicInfo{"Felix", ""}
		http.Error(w, NewValidator().Struct(info).Error(), http.StatusUnprocessableEntity)
	})

	response := httptest.NewRecorder()

	srv.ServeHTTP(response, httptest.NewRequest("GET", "/", nil))

	resp := JSONError{}
	err := json.Unmarshal(response.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf(`Unexpcted error on parse json, got :"%v" `, err.Error())
	}

}

func TestError404(t *testing.T) {

	srv := negroni.New(NewMiddleware())

	srv.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		http.Error(w, "User 420 not found", http.StatusNotFound)
	})

	response := httptest.NewRecorder()

	srv.ServeHTTP(response, httptest.NewRequest("GET", "/", nil))

	resp := JSONError{}
	err := json.Unmarshal(response.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf(`Unexpcted error on parse json, got :"%v" `, err.Error())
	}

}
