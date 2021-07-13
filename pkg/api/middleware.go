package api

import (
	"net/http"
	"time"
)

const TIMEOUT_MESSAGE = `{"error": {"code": 503,"message": "Request timeout."}}`

func timeoutMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.TimeoutHandler(next, 2*time.Second, TIMEOUT_MESSAGE).ServeHTTP(w, r)
	}
}
