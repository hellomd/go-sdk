package authentication

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// TokenExtraClaims -
type TokenExtraClaims struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	IsService *bool  `json:"isService,omitempty"`
}

type claims struct {
	TokenExtraClaims
	jwt.StandardClaims
}

type TokenCreator struct {
	secret []byte
}

// NewTokenCreator -
func NewTokenCreator(secret []byte) *TokenCreator {
	return &TokenCreator{secret}
}

func (tc TokenCreator) create(extra *TokenExtraClaims, exp time.Time) (tkn string, err error) {
	finalClaims := &claims{
		TokenExtraClaims: *extra,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp.Unix(),
		},
	}

	unsignedTkn := jwt.NewWithClaims(jwt.SigningMethodHS256, finalClaims)

	tkn, err = unsignedTkn.SignedString(tc.secret)
	return
}

func (tc TokenCreator) CreateAccessTkn(extra *TokenExtraClaims, exp time.Time) (tkn string, err error) {
	return tc.create(extra, exp)
}

func (tc TokenCreator) CreateRefreshTkn(extra *TokenExtraClaims, exp time.Time) (tkn string, err error) {
	return tc.create(extra, exp)
}
