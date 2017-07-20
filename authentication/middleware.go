package authentication

import (
	"errors"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

const headerKey = "Authorization"

var (
	errInvalidHeader = errors.New("invalid authorization header")
)

// Middleware -
type Middleware interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type middleware struct {
	secret []byte
}

// NewMiddleware -
func NewMiddleware(secret []byte) Middleware {
	return &middleware{secret}
}

func (mw *middleware) validateJWT(token *jwt.Token) (interface{}, error) {
	return mw.secret, nil
}

func (mw *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	currentUser := &CurrentUser{}

	authorization := r.Header.Get(headerKey)
	if authorization != "" {
		parts := strings.Split(authorization, " ")
		if len(parts) != 2 || parts[0] != "bearer" {
			http.Error(w, errInvalidHeader.Error(), http.StatusUnauthorized)
			return
		}

		token := parts[1]
		claims := &Claims{CurrentUser: currentUser}
		_, err := jwt.ParseWithClaims(token, claims, mw.validateJWT)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

	}
	next(w, r.WithContext(SetCurrentUserOnCtx(r.Context(), currentUser)))
}
