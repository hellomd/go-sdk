package mongo

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/urfave/negroni"
)

func TestTestMiddleware(t *testing.T) {
	db := NewTestDB()
	mw := NewTestMiddleware(db.DB)
	req := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	a := negroni.New()
	a.UseFunc(mw)
	a.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		v, err := GetFromCtx(r.Context())
		if err == errNotInCtx {
			t.Error("Expected database on context")
		}
		if v.Name != db.DBName {
			t.Errorf("Expected %s database name, got: %s", db.DBName, v.Name)
		}
		next(w, r)
	}))

	a.ServeHTTP(response, req)
}
