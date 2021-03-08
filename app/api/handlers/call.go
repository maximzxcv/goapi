package handlers

import (
	"encoding/json"
	"fmt"
	"goapi/business/data/call"
	"goapi/foundation/dbase"
	"net/http"

	"github.com/pkg/errors"
)

// CallHandler ....
type CallHandler struct {
	crep *call.CallRepository
}

// NewCallHandler constructor for CallHandler
func NewCallHandler(crep *call.CallRepository) *CallHandler {
	return &CallHandler{crep}
}

// GetCalls gets all the call for logged in user
func (ch *CallHandler) GetCalls(w http.ResponseWriter, r *http.Request) *ErrorResponse {
	ctx := r.Context()

	uid := ctx.Value("UserId")

	if uid == nil {
		return &ErrorResponse{errors.New("Cannot get userID"), http.StatusBadRequest}
	}

	cls, err := ch.crep.QueryByUser(ctx, uid)

	if err != nil {
		switch errors.Cause(err) {
		case dbase.ErrNotExist:
			return &ErrorResponse{err, http.StatusNotFound}
		default:
			return &ErrorResponse{err, 500}
		}
	}

	out, err := json.Marshal(cls)
	if err != nil {
		return &ErrorResponse{err, 500}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(out))

	return nil
}
