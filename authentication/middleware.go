package authentication

import (
	"errors"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const headerKey = "Authorization"

var errInvalidHeader = errors.New("invalid authorization header")

// NewMiddleware -
func NewMiddleware(secret []byte) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	tokenCreator := NewTokenCreator(secret)
	validateJWT := func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	}

	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		user := &User{}

		authorization := r.Header.Get(headerKey)
		if authorization != "" {
			parts := strings.Split(authorization, " ")
			if len(parts) != 2 || parts[0] != "bearer" {
				http.Error(w, errInvalidHeader.Error(), http.StatusUnauthorized)
				return
			}

			token := parts[1]
			claims := &Claims{User: user}
			_, err := jwt.ParseWithClaims(token, claims, validateJWT)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
		}

		ctx := r.Context()

		serviceToken, err := tokenCreator.CreateAccessTkn(&TokenExtraClaims{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			IsService: true,
		}, time.Now().Add(1*time.Hour))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx = SetServiceTokenInCtx(ctx, serviceToken)
		next(w, r.WithContext(SetUserInCtx(ctx, user)))
	}
}
