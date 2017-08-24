package authentication

import (
	"context"
	"errors"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	// Scheme is the Authorization header scheme
	Scheme = "bearer"

	// HeaderKey is the key used for the authentication information in headers
	HeaderKey = "Authorization"
)

var errInvalidHeader = errors.New("invalid authorization header")

// NewContextAuthenticator creates a function that stores authentication information in context
func NewContextAuthenticator(secret []byte) func(ctx context.Context, authHeader string) (context.Context, error) {
	tokenCreator := NewTokenCreator(secret)
	validateJWT := func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	}

	return func(ctx context.Context, authHeader string) (context.Context, error) {
		user := &User{}

		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != Scheme {
				return nil, errInvalidHeader
			}

			token := parts[1]
			claims := &Claims{User: user}
			_, err := jwt.ParseWithClaims(token, claims, validateJWT)
			if err != nil {
				return nil, err
			}
		}

		isService := new(bool)
		*isService = true
		serviceToken, err := tokenCreator.CreateAccessTkn(&TokenExtraClaims{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			IsService: isService,
		}, time.Now().Add(1*time.Hour))
		if err != nil {
			panic(err)
		}

		ctx = SetUserInCtx(ctx, user)
		ctx = SetServiceTokenInCtx(ctx, serviceToken)
		return ctx, nil
	}
}
