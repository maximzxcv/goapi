package handlers

import (
	"encoding/json"
	"fmt"
	"goapi/business/data/user"
	"net/http"

	"github.com/dimfeld/httptreemux"
	"github.com/pkg/errors"
)

// ErrorResponse ...
type ErrorResponse struct {
	error
	Code int
}

// UserHandler ....
type UserHandler struct {
	urep user.UserRepository
}

//NewUserHandler ...
func NewUserHandler(urep user.UserRepository) UserHandler {
	return UserHandler{
		urep: urep,
	}
}

// GetUsers ....
func (uh *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) *ErrorResponse {
	ctx := r.Context()

	uuu, err := uh.urep.Query(ctx)

	if err != nil {
		switch errors.Cause(err) {
		case user.NotExist:
			return &ErrorResponse{err, http.StatusNotFound}
		default:
			return &ErrorResponse{err, 500}
		}
	}

	out, err := json.Marshal(uuu)
	if err != nil {
		return &ErrorResponse{err, 500}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(out))

	return nil
}

// GetUserByID ....
func (uh *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) *ErrorResponse {
	ctx := r.Context()

	params := httptreemux.ContextParams(ctx)

	usr, err := uh.urep.QueryByID(ctx, params["id"])
	if err != nil {
		switch errors.Cause(err) {
		case user.NotExist:
			return &ErrorResponse{err, http.StatusNotFound}
		default:
			return &ErrorResponse{err, 500}
		}
	}
	out, err := json.Marshal(usr)
	if err != nil {
		return &ErrorResponse{err, 500}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(out))

	return nil
}

// CreateUser ....
func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) *ErrorResponse {
	ctx := r.Context()

	var nusr user.CreateUser
	if err := decode(r, &nusr); err != nil {
		return &ErrorResponse{err, 500}
	}

	usr, err := uh.urep.Create(ctx, nusr)
	if err != nil {
		return &ErrorResponse{err, 500}
	}

	out, err := json.Marshal(usr)
	if err != nil {
		return &ErrorResponse{err, 500}
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(out))

	return nil
}

// UpdateUser .....
func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) *ErrorResponse {
	ctx := r.Context()

	params := httptreemux.ContextParams(ctx)
	var uusr user.UpdateUser
	if err := decode(r, &uusr); err != nil {
		return &ErrorResponse{err, 500}
	}

	usr, err := uh.urep.Update(ctx, params["id"], uusr)
	if err != nil {
		return &ErrorResponse{err, 500}
	}

	out, err := json.Marshal(usr)
	if err != nil {
		return &ErrorResponse{err, 500}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(out))

	return nil
}

// DeleteUser .....
func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) *ErrorResponse {
	ctx := r.Context()

	params := httptreemux.ContextParams(ctx)

	if err := uh.urep.Delete(ctx, params["id"]); err != nil {
		switch errors.Cause(err) {
		case user.NotExist:
			return &ErrorResponse{err, http.StatusNotFound}
		default:
			return &ErrorResponse{err, 500}
		}
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func decode(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
