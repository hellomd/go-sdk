package recovery

// import (
// 	"bytes"
// 	"testing"

// 	"net/http"
// 	"net/http/httptest"

// 	"strings"

// 	"fmt"

// 	"reflect"

// 	raven "github.com/getsentry/raven-go"
// 	"github.com/sirupsen/logrus"
// 	"github.com/urfave/negroni"
// )

// type someTypeKey string

// const (
// 	panicMsg              = "PANIC!"
// 	somekey   someTypeKey = "someKey"
// 	someValue             = "someValue"
// )

// func TestLogRecover(t *testing.T) {
// 	errBuffer := &bytes.Buffer{}
// 	logger := logrus.New()
// 	logger.Out = errBuffer

// 	fr := &fakeRavenClient{}

// 	recMid := NewMiddleware(fr, logger)

// 	response := httptest.NewRecorder()
// 	srv := negroni.New(recMid)

// 	srv.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
// 		panic(panicMsg)
// 	})

// 	testReq := httptest.NewRequest("GET", "/", nil)
// 	srv.ServeHTTP(response, testReq)

// 	if response.Code != http.StatusInternalServerError {
// 		t.Errorf("Response Code should be %v, got %v", http.StatusInternalServerError, response.Code)
// 	}

// 	if msg := errBuffer.String(); !strings.Contains(msg, panicMsg) {
// 		t.Errorf("Logger should contains %v, but got %v", panicMsg, msg)
// 	} else {
// 		t.Log(msg)
// 	}

// 	if len(fr.errors) == 1 {
// 		if !strings.Contains(fr.errors[0], panicMsg) {
// 			t.Errorf("Raven's erros list should countains %v, got %v", panicMsg, fr.errors[0])
// 		} else {
// 			t.Log(fr.errors)
// 		}
// 	} else {
// 		t.Errorf("Size of Raven's errors list should be 1, got %v", len(fr.errors))
// 	}

// 	expHTTP := raven.NewHttp(testReq)
// 	if fr.httpContext != nil && reflect.DeepEqual(expHTTP, fr.httpContext) {
// 		t.Log(fr.httpContext)
// 	} else {
// 		t.Errorf("Raven's HTTP Context was not captured. Should be %v, got %v.", expHTTP, fr.httpContext)
// 	}

// }

// type fakeRavenClient struct {
// 	httpContext *raven.Http
// 	errors      []string
// }

// func (fr *fakeRavenClient) CaptureError(err error, tags map[string]string, interfaces ...raven.Interface) string {
// 	fr.errors = append(fr.errors, fmt.Sprintf("erros: %v; tags: %v; interfaces:%v.", err, tags, interfaces))
// 	return ""
// }

// func (fr *fakeRavenClient) SetHttpContext(http *raven.Http) {
// 	fr.httpContext = http
// }
