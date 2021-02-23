package middle

import (
	"goapi/bal"
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

// LoggMiddle   ....
func LoggMiddle(logg *bal.Logg) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					//	logg.Error(err)
				}
			}()
			start := time.Now()
			loggrw := wrapResponseWriter(w)
			next.ServeHTTP(loggrw, r)
			logg.Debug(
				"status", loggrw.status,
				"method", r.Method,
				"path", r.URL.EscapedPath(),
				"duration", time.Since(start),
			)
		}

		return http.HandlerFunc(fn)
	}
}
