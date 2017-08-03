package logger

import (
	"bytes"
	"context"
	"testing"

	"net/http/httptest"

	"strings"

	"net/http"

	"github.com/hellomd/go-sdk/requestid"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

const (
	myReqID = "42"
)

func TestBasicLogger(t *testing.T) {
	realIP := "127.0.0.1"
	errBuffer := &bytes.Buffer{}
	logger := logrus.New()
	logger.Out = errBuffer

	srv := negroni.New(NewMiddleware(logger))
	response := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil).WithContext(context.WithValue(nil, requestid.RequestIDCtxKey, myReqID))
	req.Header.Add(RealIPHeaderKey, realIP)

	srv.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.WriteHeader(http.StatusAccepted)
		next(w, r)
	})

	srv.ServeHTTP(response, req)

	if msg := errBuffer.String(); msg == "" {
		t.Error("Log error is empty")
	} else {
		if !strings.Contains(msg, myReqID) {
			t.Errorf("Expected %v in RequestID, got '%v' in log", myReqID, msg)
		}
		if !strings.Contains(msg, "status=202") {
			t.Errorf("Expected 202 in status, got '%v' in log", msg)
		}
		if !strings.Contains(msg, realIP) {
			t.Errorf("Expected %v in RealIP, got '%v' in log", realIP, msg)
		}
	}

	t.Log(errBuffer.String())
}

func TestSetsInCtx(t *testing.T) {
	realIP := "127.0.0.1"
	errBuffer := &bytes.Buffer{}
	logger := logrus.New()
	logger.Out = errBuffer

	srv := negroni.New(NewMiddleware(logger))
	response := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil).WithContext(context.WithValue(nil, requestid.RequestIDCtxKey, myReqID))
	req.Header.Add(RealIPHeaderKey, realIP)

	srv.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		entry, err := GetFromCtx(r.Context())
		if err != nil {
			t.Errorf("Expected nil error got %v", err)
		}
		if remoteEntry := entry.Data["remote"]; remoteEntry != realIP {
			t.Errorf("Expected remote %v got %v", realIP, remoteEntry)
		}
		if requestIDEntry := entry.Data["request_id"]; requestIDEntry != myReqID {
			t.Errorf("Expected request_id %v got %v", myReqID, requestIDEntry)
		}
		next(w, r)
	})

	srv.ServeHTTP(response, req)
}
