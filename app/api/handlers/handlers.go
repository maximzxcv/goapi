package handlers

import (
	"goapi/app/api/middle"
	"goapi/business/data/user"
	"log"
	"net/http"

	"github.com/dimfeld/httptreemux"
	"github.com/jmoiron/sqlx"
)

func Api(db *sqlx.DB) http.Handler {
	router := httptreemux.NewContextMux()

	urep := user.NewRepository(db)
	uh := NewUserHandler(urep)
	router.Handler(http.MethodGet, "/users/:id", appHandler(uh.GetUserByID))
	router.Handler(http.MethodPut, "/users/:id", appHandler(uh.UpdateUser))
	router.Handler(http.MethodDelete, "/users/:id", appHandler(uh.DeleteUser))
	router.Handler(http.MethodGet, "/users", appHandler(uh.GetUsers))
	router.Handler(http.MethodPost, "/users", appHandler(uh.CreateUser))

	loggMiddle := middle.LoggMiddle()
	cnfgrdRouter := loggMiddle(router)

	return cnfgrdRouter
}

// AppHandler .....
type appHandler func(http.ResponseWriter, *http.Request) *ErrorResponse

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if erresp := fn(w, r); erresp != nil {
		log.Printf("%+v", erresp)
		w.WriteHeader(erresp.Code)
	}
}
