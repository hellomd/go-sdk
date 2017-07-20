package authentication

import jwt "github.com/dgrijalva/jwt-go"

// Claims -
type Claims struct {
	*CurrentUser
	jwt.StandardClaims
}
