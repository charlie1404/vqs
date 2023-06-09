package api

import "github.com/valyala/fasthttp"

type MessageAttribute struct {
	DataType string
	Value    []byte
}

type MessageAttributes map[string]MessageAttribute

func (appCtx *AppContext) DeleteMessage(ctx *fasthttp.RequestCtx) {
}

func (appCtx *AppContext) ListQueues(ctx *fasthttp.RequestCtx) {
}

func (appCtx *AppContext) DeleteMessageBatch(ctx *fasthttp.RequestCtx) {
}

func (appCtx *AppContext) SendMessageBatch(ctx *fasthttp.RequestCtx) {
}

func (appCtx *AppContext) DeleteQueue(ctx *fasthttp.RequestCtx) {
}

func (appCtx *AppContext) TagQueue(ctx *fasthttp.RequestCtx) {
}

func (appCtx *AppContext) ListQueueTags(ctx *fasthttp.RequestCtx) {
}

func (appCtx *AppContext) UntagQueue(ctx *fasthttp.RequestCtx) {
}

func (appCtx *AppContext) SetQueueAttributes(ctx *fasthttp.RequestCtx) {
}

func (appCtx *AppContext) GetQueueAttributes(ctx *fasthttp.RequestCtx) {
}

func (appCtx *AppContext) GetQueueUrl(ctx *fasthttp.RequestCtx) {
}

func (appCtx *AppContext) PurgeQueue(ctx *fasthttp.RequestCtx) {
}

func (appCtx *AppContext) ListDeadLetterSourceQueues(ctx *fasthttp.RequestCtx) {
}

func (appCtx *AppContext) ChangeMessageVisibility(ctx *fasthttp.RequestCtx) {
}
