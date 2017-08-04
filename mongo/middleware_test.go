package mongo

import (
	"net/http"
	"net/http/httptest"
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

// We couldn't make the mongo docker image run with a ssl connection
// func TestMiddlewareWithSSL(t *testing.T) {
// 	mdw, err := NewMiddleware(config.Get(URLCfgKey)+"/test?ssl=true", true)
// 	if err != nil {
// 		t.Error("Could not create Middleware: " + err.Error())
// 	}

// 	test(mdw, t)
// }

func test(mdw Middleware, t *testing.T) {
	mw := mdw.(*middleware)
	req := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	a := negroni.New()
	a.Use(mw)

	a.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		v, err := GetFromCtx(r.Context())
		if err == errNotInCtx {
			t.Error("Expected database on context")
		}
		if v.Name != mw.dbName {
			t.Errorf("Expected %s database name, got: %s", v.Name, mw.dbName)
		}
		if v.Session == mw.session {
			t.Errorf("Expected session to be different from stored session")
		}
		next(w, r)
	}))

	a.ServeHTTP(response, req)
}
