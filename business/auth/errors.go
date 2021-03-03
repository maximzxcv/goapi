package auth

import "errors"

var (
	ErrTokenExpired    = errors.New("auth token is expied")
	ErrIncorrectHeader = errors.New("auth header not correct")
	ErrNotAuthorised   = errors.New("not authorised")
)
