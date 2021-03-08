package middle

import (
	"goapi/business/data/call"
	"log"
	"net/http"
)

// CallMiddle ....
func CallMiddle(callRepository *call.CallRepository) Middleware {
	return func(handler http.Handler) http.Handler {
		return &callmiddle{
			handler:    handler,
			repository: callRepository,
		}
	}
}

type callmiddle struct {
	handler    http.Handler
	repository *call.CallRepository
}

func (cm *callmiddle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := ctx.Value("UserId")

	ccall := call.CreateCall{
		Method: r.Method,
		Path:   r.RequestURI,
	}

	cm.handler.ServeHTTP(w, r)

	if _, err := cm.repository.CreateForUser(ctx, ccall, uid); err != nil {
		log.Print(err)
	}
}
