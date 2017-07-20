package authentication

import jwt "github.com/dgrijalva/jwt-go"

// Claims -
type Claims struct {
	*User
	jwt.StandardClaims
}
