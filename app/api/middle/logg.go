package middle

import (
	"log"
	"net/http"
	"time"
)

type loggResponseWrite struct {
	http.ResponseWriter
	status int
}

func wrapResponseWriter(w http.ResponseWriter) *loggResponseWrite {
	return &loggResponseWrite{ResponseWriter: w}
}

func (rw *loggResponseWrite) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggMiddle ....
func LoggMiddle() Middleware {
	return func(handler http.Handler) http.Handler {
		return &loggm{handler}
	}
}

type loggm struct {
	handler http.Handler
}

func (lm *loggm) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	loggrw := wrapResponseWriter(w)
	lm.handler.ServeHTTP(w, r)

	log.Println("->", loggrw.status, r.Method, r.URL.EscapedPath(), ":", time.Since(start))

}
