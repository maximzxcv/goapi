package middle

import (
	"log"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

// AuthMiddle ....
func AuthMiddle() Middleware {
	return func(handler http.Handler) http.Handler {
		return &authm{handler}
	}
}

type authm struct {
	handler http.Handler
}

func (lm *authm) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	lm.handler.ServeHTTP(w, r)

	log.Println("Auth")

}
