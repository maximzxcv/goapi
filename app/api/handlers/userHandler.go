package handlers

import (
	"encoding/json"
	"fmt"
	"goapi/bal"
	"net/http"

	"github.com/google/uuid"
)

// GetUsers ....
func GetUsers(w http.ResponseWriter, r *http.Request) {
	us := make([]bal.User, 5)
	us[0] = bal.User{
		ID:   uuid.New().String(),
		Name: "Name 1",
	}
	us[1] = bal.User{
		ID:   uuid.New().String(),
		Name: "Name 2",
	}

	out, err := json.Marshal(us)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(out))
}
