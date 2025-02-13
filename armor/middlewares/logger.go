package middlewares

import (
	"armor/utils"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggerMiddleware logs the incoming HTTP request & its duration.
func LoggerMiddleware(next http.Handler) http.Handler {
	log.Println("INFO: registering LoggerMiddleware")
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			defer func() {
				if err := recover(); err != nil {
					utils.ErrorJsonResponse(w, "something went wrong", http.StatusInternalServerError)
					log.Printf(
						"FATAL: status=500 method=%s path=%s duration=%s error=%s. Trace=%s",
						r.Method, r.URL.EscapedPath(), time.Since(start), err, debug.Stack(),
					)
				}
			}()

			wrapped := wrapResponseWriter(w)

			next.ServeHTTP(wrapped, r)

			log.Printf(
				"INFO: status=%d method=%s path=%s duration=%s", wrapped.status,
				r.Method, r.URL.EscapedPath(), time.Since(start),
			)
		},
	)
}
