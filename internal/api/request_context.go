package api

import (
	"github.com/valyala/fasthttp"
)

type FormValues map[string]string

func (f FormValues) Set(key, value string) {
	if key != "" {
		f[key] = value
	}
}

func (f FormValues) Reset() FormValues {
	for k := range f {
		f[k] = ""
	}

	return f
}

type StreamWriter struct {
	ctx *fasthttp.RequestCtx
}

func (s StreamWriter) Write(p []byte) (n int, err error) {
	s.ctx.Write(p)
	return
}

// type ExtentedRequestCtx struct {
// 	*fasthttp.RequestCtx
// 	PostForm map[string]string
// }

// type VQSRequestHandler func(ctx *ExtentedRequestCtx)

// func NewExtentedRequestCtx(ctx *fasthttp.RequestCtx) *ExtentedRequestCtx {
// 	return &ExtentedRequestCtx{
// 		RequestCtx: ctx,
// 		PostForm:   make(map[string]string),
// 	}
// }

// func (e *ExtentedRequestCtx) SetPostFormValue(key, value string) {
// 	e.PostForm[key] = value
// }
