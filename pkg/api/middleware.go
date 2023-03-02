package api

import (
	"net/http"
	"time"

	"github.com/charlie1404/vqs/pkg/o11y/logs"
	"github.com/charlie1404/vqs/pkg/utils"
)

const TIMEOUT_MESSAGE = `{"error": {"code": 503,"message": "Request timeout."}}`

func timeoutMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.TimeoutHandler(h, 2*time.Second, TIMEOUT_MESSAGE).ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func logRequestHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logId := utils.GenerateUUIDLikeId()

		logs.Logger.
			Info().
			Str("co-relation-id", logId).
			Msg("request")

		defer func() {
			logs.Logger.
				Info().
				Str("co-relation-id", logId).
				Str("method", r.Method).
				Str("url", r.URL.RequestURI()).
				Str("user_agent", r.UserAgent()).
				Dur("elapsed_ms", time.Since(start)).
				Msg("response")
		}()

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
