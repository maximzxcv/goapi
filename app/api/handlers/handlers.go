package handlers

import (
	"goapi/app/api/middle"
	"goapi/business/data/user"
	"log"
	"net/http"

	"github.com/dimfeld/httptreemux"
	"github.com/jmoiron/sqlx"
)

// ErrorResponse ...
type ErrorResponse struct {
	error
	Code int
}

// API configures URL handling
func API(db *sqlx.DB, mw ...middle.Middleware) http.Handler {
	router := apiMux{
		ContextMux: httptreemux.NewContextMux(),
		mw:         mw,
	}

	urep := user.NewRepository(db)

	uh := NewUserHandler(urep)
	router.MHandler(http.MethodGet, "/users/:id", appHandler(uh.GetUserByID), middle.AuthMiddle())
	router.MHandler(http.MethodPut, "/users/:id", appHandler(uh.UpdateUser), middle.AuthMiddle())
	router.MHandler(http.MethodDelete, "/users/:id", appHandler(uh.DeleteUser), middle.AuthMiddle())
	router.MHandler(http.MethodGet, "/users", appHandler(uh.GetUsers), middle.AuthMiddle())
	router.MHandler(http.MethodPost, "/users", appHandler(uh.CreateUser), middle.AuthMiddle())

	ah := NewAuthHandler(urep)
	router.Handler(http.MethodPost, "/singup", appHandler(ah.Singup))
	router.Handler(http.MethodPost, "/login", appHandler(ah.Login))
	router.Handler(http.MethodGet, "/logout", appHandler(ah.Logout))

	return router
}

type apiMux struct {
	*httptreemux.ContextMux
	mw []middle.Middleware
}

// MHandler allows to run middleware
func (cg *apiMux) MHandler(method, path string, handler http.Handler, m ...middle.Middleware) {
	h := wrapMiddleware(cg.mw, handler)
	h = wrapMiddleware(m, h)
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
