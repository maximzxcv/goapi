package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Signup ...
type Signup struct {
	Username string
	Password string
}

// Login ...
type Login struct {
	Username string
	Password string
}

// Access ...
type Access struct {
	Token string
}

// ApiClaim ...
type ApiClaim struct {
	jwt.StandardClaims
	Username string
}

func NewClaim(username string) ApiClaim {
	return ApiClaim{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
			Issuer:    "goapi",
		},
	}
}
