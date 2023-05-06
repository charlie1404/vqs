package api

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/charlie1404/vqs/pkg/o11y/metrics"
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
		// logId := utils.GenerateUUIDLikeId()

		// logs.Logger.
		// 	Info().
		// 	Str("co-relation-id", logId).
		// 	Msg("request")

		// defer func() {
		// 	logs.Logger.
		// 		Info().
		// 		Str("co-relation-id", logId).
		// 		Msg("response")
		// }()

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

type HTTPMetrics struct {
	http.ResponseWriter
	statusCode int
	written    int
}

var httpMetricsPool = sync.Pool{New: func() interface{} { return &HTTPMetrics{} }}

func (hm *HTTPMetrics) WriteHeader(code int) {
	hm.statusCode = code
	hm.ResponseWriter.WriteHeader(code)
}

func validatePostRequestAndAction(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			resp := toXMLErrorResponse("UnsupportedMethod", fmt.Sprintf("(%s) method is not supported", r.Method), "Kuch bhi")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(resp)
			return
		}

		r.ParseForm()

		if r.Form.Get("Action") == "" {
			resp := toXMLErrorResponse("MissingAction", "The request must contain the parameter Action.", "")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(resp)
			return
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func httpMetricsMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		m := httpMetricsPool.Get().(*HTTPMetrics)
		start := time.Now()

		m.ResponseWriter = w
		m.statusCode = http.StatusOK
		m.written = 0

		defer func() {
			defer httpMetricsPool.Put(m)
			statusText := http.StatusText(m.statusCode)
			metrics.IncHttpRequestsCounter(statusText, r.FormValue("Action"))
			metrics.ObserveHttpRequestsDuration(statusText, r.FormValue("Action"), float64(time.Since(start).Microseconds()))
		}()

		h.ServeHTTP(m, r)
	}

	return http.HandlerFunc(fn)
}

func Middleware(h http.Handler) http.Handler {
	return logRequestHandler(
		timeoutMiddleware(
			validatePostRequestAndAction(
				httpMetricsMiddleware(
					h,
				),
			),
		),
	)
}
