package errors

import (
	"bytes"
	"errors"
	"testing"

	"net/http"

	"net/http/httptest"

	"strings"

	"encoding/json"

	"fmt"

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

	if response.Body.String() != http.StatusText(http.StatusInternalServerError) {
		t.Errorf(`Unexpcted body, got :"%v" `, response.Body.String())
	}

	if !strings.Contains(errBuffer.String(), "chaos") {
		t.Errorf("Unexpcted log, got :%v ", errBuffer.String())
	}

}

func TestError400(t *testing.T) {
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

	resp := []inError{}
	err := json.Unmarshal(response.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf(`Unexpcted error on parse json, got :"%v" `, err.Error())
	}

	fmt.Println(resp ) 

}
