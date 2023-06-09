package api

import (
	"bytes"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/charlie1404/vqs/pkg/o11y/metrics"
	"github.com/charlie1404/vqs/pkg/utils"
)

var (
	TIMEOUT_MESSAGE    = `{"error": {"code": 503,"message": "Request timeout."}}`
	QUERY_PARAM_SEP    = byte('&')
	QUERY_PARAM_KV_SEP = byte('=')
)

func setXmlContentType(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	fn := func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/xml")
		h(ctx)
	}

	return fn
}

func timeoutHandler(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.TimeoutWithCodeHandler(h, 2*time.Second, TIMEOUT_MESSAGE, fasthttp.StatusServiceUnavailable)
}

func httpMetricsHandler(h fasthttp.RequestHandler) fasthttp.RequestHandler {
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

func validateContentType(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	fn := func(ctx *fasthttp.RequestCtx) {
		if string(ctx.Request.Header.ContentType()) != "application/x-www-form-urlencoded" {
			resp := toXMLErrorResponse("InvalidContentType", "Invalid content type", "")
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.SetBody(resp)
			return
		}

		h(ctx)
	}

	return fn
}

func parseHttpBodyMiddleware(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	fn := func(ctx *fasthttp.RequestCtx) {
		urlEncodedBody := ctx.PostBody()
		formValues := make(FormValues)

		for {
			pos := bytes.IndexByte(urlEncodedBody, QUERY_PARAM_SEP)
			if pos < 0 {
				break
			}

			key, value := utils.ParseUrlEncodedBodyParamKV(urlEncodedBody[:pos], QUERY_PARAM_KV_SEP)
			formValues[key] = value

			urlEncodedBody = urlEncodedBody[pos+1:]
		}

		key, value := utils.ParseUrlEncodedBodyParamKV(urlEncodedBody, QUERY_PARAM_KV_SEP)
		formValues[key] = value

		ctx.SetUserValue("body", formValues)

		h(ctx)
	}

	return fn
}

func validateRequestBasic(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	fn := func(ctx *fasthttp.RequestCtx) {
		formValues := ctx.UserValue("body").(FormValues)

		method := ctx.Method()
		if string(method) != fasthttp.MethodPost {
			resp := toXMLErrorResponse("UnsupportedMethod", fmt.Sprintf("(%s) method is not supported", method), "Kuch bhi")
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			ctx.SetBody(resp)
			return
		}

		action := formValues["Action"]
		if action == "" {
			resp := toXMLErrorResponse("MissingAction", "The request must contain the parameter Action.", "")
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.SetBody(resp)
			return
		}

		ver := formValues["Version"]
		if ver != "2012-11-05" {
			resp := toXMLErrorResponse("NoSuchVersion", fmt.Sprintf("The requested version ( %s ) is not valid.", ver), "")
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.Write(resp)
			return
		}

		h(ctx)
	}

	return fn
}

func logRequestHandler(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	fn := func(ctx *fasthttp.RequestCtx) {
		// logId := utils.GenerateUUIDLikeId()

		// logs.Logger.
		// 	Info().
		// 	Str("co-relation-id", logId).
		// 	Msg("request")

		// // defer func() {
		// // 	logs.Logger.
		// // 		Info().
		// // 		Str("co-relation-id", logId).
		// // 		Msg("response")
		// // }()

		h(ctx)
	}

	return fn
}

func Middleware(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return setXmlContentType(
		timeoutHandler(
			httpMetricsHandler(
				validateContentType(
					parseHttpBodyMiddleware(
						validateRequestBasic(
							logRequestHandler(
								h,
							),
						),
					),
				),
			),
		),
	)
}
