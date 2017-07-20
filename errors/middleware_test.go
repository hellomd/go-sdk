package errors

import (
	"bytes"
	"errors"
	"testing"

	"net/http"

	"net/http/httptest"

	"strings"

	"encoding/json"

	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

func TestError500(t *testing.T) {
	errBuffer := &bytes.Buffer{}
	logger := logrus.New()
	logger.Out = errBuffer

	srv := negroni.New(NewMiddleware(logger))

	srv.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		http.Error(w, errors.New("chaos").Error(), http.StatusInternalServerError)
	})

	response := httptest.NewRecorder()

	srv.ServeHTTP(response, httptest.NewRequest("GET", "/", nil))

	resp := &JSONError{}
	err := json.NewDecoder(response.Body).Decode(resp)
	if err != nil {
		t.Errorf(`Unexpcted error on parse json, got :"%v" `, err.Error())
	}

	if !reflect.DeepEqual(resp, ErrUnexptectedError) {
		t.Errorf(`Unexpcted error response. Got :"%v", want: %v `, resp, ErrUnexptectedError)
	}

	if !strings.Contains(errBuffer.String(), "chaos") {
		t.Errorf("Unexpcted log, got :%v ", errBuffer.String())
	}

}

func TestError422(t *testing.T) {
	errBuffer := &bytes.Buffer{}
	logger := logrus.New()
	logger.Out = errBuffer

	srv := negroni.New(NewMiddleware(logger))

	srv.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		info := BasicInfo{"Felix", ""}
		http.Error(w, validate.Struct(info).Error(), http.StatusUnprocessableEntity)
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
	errBuffer := &bytes.Buffer{}
	logger := logrus.New()
	logger.Out = errBuffer

	srv := negroni.New(NewMiddleware(logger))

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
