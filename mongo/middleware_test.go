package mongo

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hellomd/go-sdk/config"
	"github.com/urfave/negroni"
)

func TestMiddlewareWithoutSSL(t *testing.T) {
	mdw, err := NewMiddleware(config.Get(URLCfgKey)+"/test", false)
	if err != nil {
		t.Error("Could not create Middleware: " + err.Error())
	}

	test(mdw, t)
}

func TestMiddlewareWithSSL(t *testing.T) {
	mdw, err := NewMiddleware(config.Get(URLCfgKey)+"/test?ssl=true", true)
	if err != nil {
		t.Error("Could not create Middleware: " + err.Error())
	}

	test(mdw, t)
}

func test(mdw Middleware, t *testing.T) {
	mw := mdw.(*middleware)
	req := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	a := negroni.New()
	a.Use(mw)

	a.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		dbName := strings.TrimPrefix(mw.mongoURL.Path, "/")
		v, err := GetFromCtx(r.Context())
		if err == errNotInCtx {
			t.Error("Expected database on context")
		}
		if v.Name != dbName {
			t.Errorf("Expected %s database name, got: %s", v.Name, dbName)
		}
		if v.Session == mw.session {
			t.Errorf("Expected session to be different from stored session")
		}
		next(w, r)
	}))

	a.ServeHTTP(response, req)
}
