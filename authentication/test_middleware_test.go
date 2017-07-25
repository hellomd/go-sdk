package authentication

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/urfave/negroni"
)

func TestTestMiddlewareDefault(t *testing.T) {
	response := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	srv := negroni.New()
	srv.Use(NewTestMiddleware())
	srv.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		user := GetUserFromCtx(r.Context())
		if !user.Empty() {
			t.Errorf("user should be empty")
		}
	}))
	srv.ServeHTTP(response, req)
}

func TestTestMiddlewareOverwrite(t *testing.T) {
	response := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	expectedUser := &User{ID: "someone"}
	middleware := NewTestMiddleware()
	middleware.SetUser(expectedUser)

	srv := negroni.New()
	srv.Use(middleware)
	srv.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		user := GetUserFromCtx(r.Context())
		if user != expectedUser {
			t.Errorf("expected user to be %v, got %v", expectedUser, user)
		}
	}))
	srv.ServeHTTP(response, req)
}
