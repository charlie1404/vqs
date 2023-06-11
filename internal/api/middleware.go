package api

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/charlie1404/vqs/internal/o11y/metrics"
	"github.com/charlie1404/vqs/internal/utils"
)

var (
	TIMEOUT_MESSAGE    = `{"error": {"code": 503,"message": "Request timeout."}}`
	QUERY_PARAM_SEP    = byte('&')
	QUERY_PARAM_KV_SEP = byte('=')
)

type Middleware struct {
	templates *template.Template
}

func (m *Middleware) timeoutHandler(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.TimeoutWithCodeHandler(h, 2*time.Second, TIMEOUT_MESSAGE, fasthttp.StatusServiceUnavailable)
}

func (m *Middleware) setXmlContentType(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	fn := func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/xml")
		h(ctx)
	}

	return fn
}

func (m *Middleware) httpMetricsHandler(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	fn := func(ctx *fasthttp.RequestCtx) {
		start := time.Now()

		defer func() {
			elapsedMs := float64(time.Since(start).Microseconds())
			action := string(ctx.FormValue("Action"))

			statusText := fmt.Sprintf("%d", ctx.Response.StatusCode())
			metrics.IncHttpRequestsCounter(statusText, action)
			metrics.ObserveHttpRequestsDuration(statusText, action, elapsedMs)
		}()

		h(ctx)
	}

	return fn
}

func (m *Middleware) validateContentType(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	fn := func(ctx *fasthttp.RequestCtx) {
		x := StreamWriter{ctx}

		if string(ctx.Request.Header.ContentType()) != "application/x-www-form-urlencoded" {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			m.templates.ExecuteTemplate(x, "error.tpl", NewVQSError("InvalidContentType", "Invalid content type", ""))
			return
		}

		h(ctx)
	}

	return fn
}

func (m *Middleware) parseHttpBodyMiddleware(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	var FormValuesSyncPool = sync.Pool{
		New: func() interface{} { return make(FormValues) },
	}

	fn := func(ctx *fasthttp.RequestCtx) {
		formValues := FormValuesSyncPool.Get().(FormValues).Reset()
		defer FormValuesSyncPool.Put(formValues)

		urlEncodedBody := ctx.PostBody()

		for {
			pos := bytes.IndexByte(urlEncodedBody, QUERY_PARAM_SEP)
			if pos < 0 {
				break
			}

			key, value := utils.ParseUrlEncodedBodyParamKV(urlEncodedBody[:pos], QUERY_PARAM_KV_SEP)
			formValues.Set(key, value)

			urlEncodedBody = urlEncodedBody[pos+1:]
		}

		key, value := utils.ParseUrlEncodedBodyParamKV(urlEncodedBody, QUERY_PARAM_KV_SEP)
		formValues.Set(key, value)

		ctx.SetUserValue("body", formValues)

		h(ctx)
	}

	return fn
}

func (m *Middleware) validateRequestBasic(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	fn := func(ctx *fasthttp.RequestCtx) {
		formValues := ctx.UserValue("body").(FormValues)

		x := StreamWriter{ctx}

		method := ctx.Method()
		if string(method) != fasthttp.MethodPost {
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			m.templates.ExecuteTemplate(x, "error.tpl", NewVQSError("UnsupportedMethod", fmt.Sprintf("(%s) method is not supported", method), "Kuch bhi"))
			return
		}

		action := formValues["Action"]
		if action == "" {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			m.templates.ExecuteTemplate(x, "error.tpl", NewVQSError("MissingAction", "The request must contain the parameter Action.", ""))
			return
		}

		ver := formValues["Version"]
		if ver != "2012-11-05" {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			m.templates.ExecuteTemplate(x, "error.tpl", NewVQSError("NoSuchVersion", fmt.Sprintf("The requested version ( %s ) is not valid.", ver), ""))
			return
		}

		h(ctx)
	}

	return fn
}

func (m *Middleware) logRequestHandler(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	fn := func(ctx *fasthttp.RequestCtx) {
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

		h(ctx)
	}

	return fn
}

func (m *Middleware) WrapHandler(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return m.timeoutHandler(
		m.setXmlContentType(
			m.httpMetricsHandler(
				m.validateContentType(
					m.parseHttpBodyMiddleware(
						m.validateRequestBasic(
							m.logRequestHandler(
								h,
							),
						),
					),
				),
			),
		),
	)
}
