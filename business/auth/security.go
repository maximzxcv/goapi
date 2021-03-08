package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	salt = "8dhoer8gjodfger4"
)

// SignAccess creates access token
func (claim *ApiClaim) SignAccess() (*Access, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString([]byte(salt))
	if err != nil {
		return nil, err
	}

	return &Access{Token: signedToken}, nil
}

// ValidateAccess checks if auth token is valid
func ValidateAccess(authString string) (*ApiClaim, error) {
	parts := strings.Split(authString, " ")

	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, ErrIncorrectHeader
	}

	token, err := jwt.ParseWithClaims(
		parts[1],
		&ApiClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(salt), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*ApiClaim)

	if !ok {
		return nil, errors.New("couldn't parse claims")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return nil, ErrTokenExpired
	} else {
		return claims, nil
	}
}
