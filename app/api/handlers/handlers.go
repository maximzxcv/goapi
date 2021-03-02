package handlers

import (
	"goapi/app/api/middle"
	"goapi/business/data/user"
	"log"
	"net/http"

	"github.com/dimfeld/httptreemux"
	"github.com/jmoiron/sqlx"
)

// Api configures URL handling
func Api(db *sqlx.DB, m ...middle.Middleware) http.Handler {
	//router := httptreemux.NewContextMux()

	router := apiMux{
		httptreemux.NewContextMux(),
	}

	urep := user.NewRepository(db)
	uh := NewUserHandler(urep)

	router.MHandler(http.MethodGet, "/users/:id", appHandler(uh.GetUserByID), middle.AuthMiddle())

	router.MHandler(http.MethodPut, "/users/:id", appHandler(uh.UpdateUser), middle.AuthMiddle())
	router.MHandler(http.MethodDelete, "/users/:id", appHandler(uh.DeleteUser), middle.AuthMiddle())
	router.MHandler(http.MethodGet, "/users", appHandler(uh.GetUsers), middle.AuthMiddle())
	router.MHandler(http.MethodPost, "/users", appHandler(uh.CreateUser), middle.AuthMiddle())

	router.Handler(http.MethodPost, "/login", appHandler(uh.GetUserByID))
	router.Handler(http.MethodPost, "/logout", appHandler(uh.GetUserByID))

	h := wrapMiddleware(m, router)

	return h
}

type apiMux struct {
	*httptreemux.ContextMux
}

// MHandler allows to run middleware
func (cg *apiMux) MHandler(method, path string, handler http.Handler, m ...middle.Middleware) {
	h := wrapMiddleware(m, handler)
	cg.Handler(method, path, h)
}

// AppHandler .....
type appHandler func(http.ResponseWriter, *http.Request) *ErrorResponse

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if erresp := fn(w, r); erresp != nil {
		log.Printf("%+v", erresp)
		w.WriteHeader(erresp.Code)
	}
}

func wrapMiddleware(mw []middle.Middleware, handler http.Handler) http.Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}
