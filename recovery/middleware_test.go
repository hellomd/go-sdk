package recovery

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hellomd/go-sdk/recovery/sentry"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/urfave/negroni"
)

func TestRecoveryMiddleware(t *testing.T) {
	Convey("RecoveryMiddleware", t, func() {
		var (
			receivedErr error
			receivedCtx context.Context
		)

		handleError := func(ctx context.Context, err error) {
			receivedErr = err
			receivedCtx = ctx
		}

		response := httptest.NewRecorder()
		pipeline := negroni.New()
		pipeline.Use(&Middleware{"", handleError})

		request := func() {
			pipeline.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))
		}

		Convey("it injects sentry into context", func() {
			pipeline.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
				sentry, err := sentry.GetFromCtx(r.Context())
				So(err, ShouldBeNil)
				So(sentry, ShouldNotBeNil)
			})

			request()
		})

		Convey("it recovers on panic", func() {
			expectedErr := errors.New("test error")

			pipeline.UseFunc(func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
				panic(expectedErr)
			})

			request()

			So(receivedErr.Error(), ShouldEqual, expectedErr.Error())
			So(response.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})
}
