package handlers

import (
	"goapi/app/api/middle"
	"goapi/business/data/call"
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
		mware:      make([]middle.Middleware, 2), //  mw,
	}

	urep := user.NewRepository(db)
	crep := call.NewRepository(db)

	//	Handler: handlers.API(db, middle.LoggMiddle(), middle.CallMiddle()),
	router.mware[0] = middle.LoggMiddle()
	router.mware[1] = middle.CallMiddle(crep)

	uh := NewUserHandler(urep)
	router.MHandler(http.MethodGet, "/users/:id", appHandler(uh.GetUserByID), middle.AuthMiddle())
	router.MHandler(http.MethodPut, "/users/:id", appHandler(uh.UpdateUser), middle.AuthMiddle())
	router.MHandler(http.MethodDelete, "/users/:id", appHandler(uh.DeleteUser), middle.AuthMiddle())
	router.MHandler(http.MethodGet, "/users", appHandler(uh.GetUsers), middle.AuthMiddle())
	router.MHandler(http.MethodPost, "/users", appHandler(uh.CreateUser), middle.AuthMiddle())

	ah := NewAuthHandler(urep)
	router.MHandler(http.MethodPost, "/singup", appHandler(ah.Singup), nil)
	router.MHandler(http.MethodPost, "/login", appHandler(ah.Login), nil)
	router.MHandler(http.MethodGet, "/logout", appHandler(ah.Logout), nil)

	ch := NewCallHandler(crep)
	router.MHandler(http.MethodGet, "/calls", appHandler(ch.GetCalls), middle.AuthMiddle())

	return router
}

type apiMux struct {
	*httptreemux.ContextMux
	mware []middle.Middleware
}

// MHandler allows to run middleware
func (cg *apiMux) MHandler(method, path string, handler http.Handler, m ...middle.Middleware) {
	h := wrapMiddleware(cg.mware, handler)
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
