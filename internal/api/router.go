package api

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func (appCtx *AppContext) requestHandler(ctx *fasthttp.RequestCtx) {
	formValues := ctx.UserValue("body").(FormValues)
	action := formValues["Action"]

	switch action {
	case "ListQueues":
		appCtx.ListQueues(ctx)
	case "CreateQueue":
		appCtx.CreateQueue(ctx)
	case "SendMessage":
		appCtx.SendMessage(ctx)
	case "ReceiveMessage":
		appCtx.ReceiveMessage(ctx)
	case "DeleteMessage":
		appCtx.DeleteMessage(ctx)
	case "DeleteMessageBatch":
		appCtx.DeleteMessageBatch(ctx)
	case "SendMessageBatch":
		appCtx.SendMessageBatch(ctx)
	case "DeleteQueue":
		appCtx.DeleteQueue(ctx)
	case "TagQueue":
		appCtx.TagQueue(ctx)
	case "ListQueueTags":
		appCtx.ListQueueTags(ctx)
	case "UntagQueue":
		appCtx.UntagQueue(ctx)
	case "SetQueueAttributes":
		appCtx.SetQueueAttributes(ctx)
	case "GetQueueAttributes":
		appCtx.GetQueueAttributes(ctx)
	case "GetQueueUrl":
		appCtx.GetQueueUrl(ctx)
	case "PurgeQueue":
		appCtx.PurgeQueue(ctx)
	case "ListDeadLetterSourceQueues":
		appCtx.ListDeadLetterSourceQueues(ctx)
	case "ChangeMessageVisibility":
		appCtx.ChangeMessageVisibility(ctx)
	default:
		x := StreamWriter{ctx}
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		appCtx.templates.ExecuteTemplate(x, "error.xml", NewVQSError("InvalidAction", fmt.Sprintf("The action %s is not valid for this endpoint.", action), ""))
	}
}
