package mongo

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hellomd/go-sdk/config"
	"github.com/urfave/negroni"

	"gopkg.in/mgo.v2"
)

func TestMiddleware(t *testing.T) {
	dbName := "test"
	mw := NewMiddleware(config.Get(URLCfgKey), dbName).(*middleware)
	req := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	a := negroni.New()
	a.Use(mw)
	a.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		v, ok := r.Context().Value(ctxKey{}).(*mgo.Database)
		if !ok {
			t.Error("Expected database on context")
		}
		if v.Name != dbName {
			t.Errorf("Expected %s database name, got: %s", dbName, v.Name)
		}
		if v.Session == mw.session {
			t.Errorf("Expected session to be different from stored session")
		}
		next(w, r)
	}))

	a.ServeHTTP(response, req)
}
