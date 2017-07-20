package authentication

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/urfave/negroni"
)

var secret = []byte("123456")

func newToken(userID string) string {
	claims := Claims{
		&CurrentUser{ID: userID},
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

func TestSomething(t *testing.T) {
	srv := negroni.New(NewMiddleware(secret))
	response := httptest.NewRecorder()
	fmt.Println(newToken("123"))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add(headerKey, "bearer "+newToken("123"))

	srv.ServeHTTP(response, req)
}
