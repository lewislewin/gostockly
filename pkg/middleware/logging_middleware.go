package middleware

import (
	"net/http"
	"time"

	"gostockly/pkg/logger"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Log the incoming request
		log := logger.GetLogger()
		log.Info("Incoming request: method=%s, path=%s, remote_addr=%s",
			r.Method, r.URL.Path, r.RemoteAddr)

		// Process the next handler
		next.ServeHTTP(w, r)

		// Log the duration of the request
		duration := time.Since(startTime)
		log.Info("Completed request: method=%s, path=%s, duration=%s",
			r.Method, r.URL.Path, duration)
	})
}
