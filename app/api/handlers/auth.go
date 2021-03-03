package handlers

import (
	"encoding/json"
	"fmt"
	"goapi/business/auth"
	"goapi/business/data/user"
	"goapi/foundation/dbase"
	"net/http"

	"github.com/pkg/errors"
)

// AuthHandler ....
type AuthHandler struct {
	urep *user.UserRepository
}

// NewAuthHandler constructor for AuthHandler
func NewAuthHandler(urep *user.UserRepository) *AuthHandler {
	return &AuthHandler{urep}
}

// Singup ....
func (ah *AuthHandler) Singup(w http.ResponseWriter, r *http.Request) *ErrorResponse {
	ctx := r.Context()

	var signup auth.Signup
	if err := decode(r, &signup); err != nil {
		return &ErrorResponse{err, 500}
	}
	nusr := user.CreateUser{
		Name:     signup.Username,
		Password: signup.Password,
	}
	if _, err := ah.urep.Create(ctx, nusr); err != nil {

		switch errors.Cause(err) {
		case dbase.ErrAlreadyExist:
			return &ErrorResponse{err, http.StatusConflict}
		default:
			return &ErrorResponse{err, 500}
		}

	}

	w.WriteHeader(http.StatusCreated)
	return nil
}

// Login ....
func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) *ErrorResponse {
	ctx := r.Context()

	var login auth.Login
	if err := decode(r, &login); err != nil {
		return &ErrorResponse{err, 500}
	}

	if err := ah.urep.CheckAuth(ctx, login.Username, login.Password); err != nil {
		switch errors.Cause(err) {
		case auth.ErrNotAuthorised:
			w.WriteHeader(http.StatusUnauthorized)
			return nil
		default:
			return &ErrorResponse{err, 500}
		}
	}

	claim := auth.NewClaim(login.Username)
	access, err := claim.SignAccess()
	if err != nil {
		return &ErrorResponse{err, 500}
	}

	out, err := json.Marshal(&access)
	if err != nil {
		return &ErrorResponse{err, 500}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(out))
	return nil

}

// Logout ....
func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) *ErrorResponse {
	return nil
}
