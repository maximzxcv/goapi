package handlers

import (
	"encoding/json"
	"fmt"
	"goapi/business/data/user"
	"net/http"
)

type UserHandler struct {
	urep user.UserRepository
}

func NewUserHandler(urep user.UserRepository) UserHandler {
	return UserHandler{
		urep: urep,
	}
}

// GetUsers ....
func (uh *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// us := make([]user.User, 5)
	// us[0] = user.User{
	// 	ID:   uuid.New().String(),
	// 	Name: "Name 1",
	// }
	// us[1] = user.User{
	// 	ID:   uuid.New().String(),
	// 	Name: "Name 2",
	// }
	ctx := r.Context()
	uuu, err := uh.urep.Query(ctx)

	out, err := json.Marshal(uuu)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(out))
}
