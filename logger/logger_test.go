package logger

import (
	"bytes"
	"context"
	"testing"

	"net/http/httptest"

	"strings"

	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

const (
	myReqID = "42"
)

func TestBasicLogger(t *testing.T) {
	//Prepare Logger
	errBuffer := &bytes.Buffer{}
	logger := logrus.New()
	logger.Out = errBuffer

	//Prepare server, response and request
	srv := negroni.New(NewLogger(logger))
	response := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil).WithContext(context.WithValue(nil, RequestIDcontextKey, myReqID))

	//Set handler to set StatusOK header
	srv.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.WriteHeader(http.StatusAccepted)
		next(w, r)
	})

	srv.ServeHTTP(response, req)

	//Assertions
	if msg := errBuffer.String(); msg == "" {
		t.Error("Log error is empty")
	} else {
		if !strings.Contains(msg, myReqID) {
			t.Errorf("Expected %v in RequestID, got '%v' in log", myReqID, msg)
		}
		if !strings.Contains(msg, "status=202") {
			t.Errorf("Expected 202 in status, got '%v' in log", msg)
		}
	}

	t.Log(errBuffer.String())
}
