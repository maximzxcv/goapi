package middle

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"goapi/business/auth"
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
	authStr := r.Header.Get("authorization")
	
	if claims, err := auth.ValidateAccess(authStr); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, err.Error())
	} else {
		ctx := context.WithValue(r.Context(), "Username", claims.Username)
		lm.handler.ServeHTTP(w, r.WithContext(ctx))
	}
}
