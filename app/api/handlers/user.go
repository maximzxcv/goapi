package handlers

import (
	"encoding/json"
	"fmt"
	"goapi/business/data/user"
	"log"
	"net/http"

	"github.com/dimfeld/httptreemux"
)

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
func (uh *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	uuu, err := uh.urep.Query(ctx)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), 500)
		return
	}

	out, err := json.Marshal(uuu)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(out))
}

// GetUserByID ....
func (uh *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := httptreemux.ContextParams(ctx)

	usr, err := uh.urep.QueryByID(ctx, params["id"])
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), 500)
		return
	}

	out, err := json.Marshal(usr)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(out))
}

// CreateUser ....
func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var nusr user.CreateUser
	if err := decode(r, &nusr); err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), 500)
		return
	}

	usr, err := uh.urep.Create(ctx, nusr)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), 500)
		return
	}

	out, err := json.Marshal(usr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(out))
}

// UpdateUser .....
func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := httptreemux.ContextParams(ctx)
	var uusr user.UpdateUser
	if err := decode(r, &uusr); err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), 500)
		return
	}

	usr, err := uh.urep.Update(ctx, params["id"], uusr)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), 500)
		return
	}

	out, err := json.Marshal(usr)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(out))
}

// DeleteUser .....
func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := httptreemux.ContextParams(ctx)

	if err := uh.urep.Delete(ctx, params["id"]); err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func decode(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return err
	}

	return nil
}
