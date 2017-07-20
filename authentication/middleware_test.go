package authentication

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/urfave/negroni"
)

var secret = []byte("123456")

func newToken(userID string) string {
	claims := Claims{
		&User{ID: userID},
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(10 * time.Hour).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		panic(err)
	}
	return token
}

func TestNoToken(t *testing.T) {
	srv := negroni.New()
	srv.UseFunc(NewMiddleware(secret))
	response := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	srv.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		user := GetUserFromCtx(r.Context())
		if !user.Empty() {
			t.Errorf("user should be empty")
		}
	}))

	srv.ServeHTTP(response, req)
}

func TestBadAuthorization(t *testing.T) {
	srv := negroni.New()
	srv.UseFunc(NewMiddleware(secret))

	badHeaders := []string{
		"something else",
		"something",
		"bearer ",
		"bearer abc",
		"bearer abc def",
	}

	for _, header := range badHeaders {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Add(headerKey, header)
		srv.ServeHTTP(response, req)

		if response.Code != http.StatusUnauthorized {
			t.Errorf("expected status code to be %v, got %v for header %v", http.StatusUnauthorized, response.Code, header)
		}
	}
}

func TestGoodToken(t *testing.T) {
	userID := "4d88e15b60f486e428412dc7"
	srv := negroni.New()
	srv.UseFunc(NewMiddleware(secret))
	srv.Use(negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		user := GetUserFromCtx(r.Context())
		if user.Empty() {
			t.Errorf("user should not be empty")
		}
		if user.ID != userID {
			t.Errorf("expected user id %v, got %v", userID, user.ID)
		}
	}))

	response := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add(headerKey, "bearer "+newToken(userID))
	srv.ServeHTTP(response, req)

	if response.Code != http.StatusOK {
		t.Errorf("expected status code to be %v, got %v", http.StatusOK, response.Code)
	}
}
