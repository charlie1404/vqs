package api

import (
	"fmt"
	"log"
	"net/http"
)

func (appCtx *AppContext) ListQueues(w http.ResponseWriter, r *http.Request) {

}

func (appCtx *AppContext) CreateQueue(w http.ResponseWriter, r *http.Request) {
	createQueueInput := parseCreateQueueInput(r.Form)

	if err := appCtx.validator.validateCreateQueueInput(createQueueInput); err != nil {
		resp := toXMLErrorResponse("InvalidAttributeValue", "Invalid value for some parameter.", "")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/xml")
		w.Write(resp)
		return
	}

	if err := appCtx.queues.CreateQueue(createQueueInput); err != nil {
		log.Println(err)
		resp := toXMLErrorResponse("InternalServerError", "Todo return better errors", "")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/xml")
		w.Write(resp)
		return
	}

	resp := toCreateQueueResponse(createQueueInput.QueueName)
	w.Header().Set("Content-Type", "application/xml")
	w.Write(resp)
}

func (appCtx *AppContext) SendMessage(w http.ResponseWriter, r *http.Request) {

}

func (appCtx *AppContext) ReceiveMessage(w http.ResponseWriter, r *http.Request) {

}

func (appCtx *AppContext) DeleteMessage(w http.ResponseWriter, r *http.Request) {

}

func (appCtx *AppContext) DeleteMessageBatch(w http.ResponseWriter, r *http.Request) {

}

func (appCtx *AppContext) SendMessageBatch(w http.ResponseWriter, r *http.Request) {

}

func (appCtx *AppContext) DeleteQueue(w http.ResponseWriter, r *http.Request) {

}

func (appCtx *AppContext) TagQueue(w http.ResponseWriter, r *http.Request) {

}

func (appCtx *AppContext) ListQueueTags(w http.ResponseWriter, r *http.Request) {

}

func (appCtx *AppContext) UntagQueue(w http.ResponseWriter, r *http.Request) {

}

func (appCtx *AppContext) SetQueueAttributes(w http.ResponseWriter, r *http.Request) {

}

func (appCtx *AppContext) GetQueueAttributes(w http.ResponseWriter, r *http.Request) {

}

func (appCtx *AppContext) GetQueueUrl(w http.ResponseWriter, r *http.Request) {

}

func (appCtx *AppContext) PurgeQueue(w http.ResponseWriter, r *http.Request) {

}

func (appCtx *AppContext) ListDeadLetterSourceQueues(w http.ResponseWriter, r *http.Request) {

}

func (appCtx *AppContext) ChangeMessageVisibility(w http.ResponseWriter, r *http.Request) {

}

func (ctx *AppContext) requestHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.Method != "POST" {
		resp := toXMLErrorResponse("UnsupportedMethod", fmt.Sprintf("(%s) method is not supported", r.Method), "")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/xml")
		w.Write(resp)
		return
	}

	if r.FormValue("Version") != "2012-11-05" {
		resp := toXMLErrorResponse("NoSuchVersion", fmt.Sprintf("The requested version ( %s ) is not valid.", r.FormValue("Version")), "")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/xml")
		w.Write(resp)
		return
	}

	switch r.FormValue("Action") {
	case "ListQueues":
		{
			ctx.ListQueues(w, r)
			break
		}
	case "CreateQueue":
		{
			ctx.CreateQueue(w, r)
			break
		}
	case "SendMessage":
		{
			ctx.SendMessage(w, r)
			break
		}
	case "ReceiveMessage":
		{
			ctx.ReceiveMessage(w, r)
			break
		}
	case "DeleteMessage":
		{
			ctx.DeleteMessage(w, r)
			break
		}
	case "DeleteMessageBatch":
		{
			ctx.DeleteMessageBatch(w, r)
			break
		}
	case "SendMessageBatch":
		{
			ctx.SendMessageBatch(w, r)
			break
		}
	case "DeleteQueue":
		{
			ctx.DeleteQueue(w, r)
			break
		}
	case "TagQueue":
		{
			ctx.TagQueue(w, r)
			break
		}
	case "ListQueueTags":
		{
			ctx.ListQueueTags(w, r)
			break
		}
	case "UntagQueue":
		{
			ctx.UntagQueue(w, r)
			break
		}
	case "SetQueueAttributes":
		{
			ctx.SetQueueAttributes(w, r)
			break
		}
	case "GetQueueAttributes":
		{
			ctx.GetQueueAttributes(w, r)
			break
		}
	case "GetQueueUrl":
		{
			ctx.GetQueueUrl(w, r)
			break
		}
	case "PurgeQueue":
		{
			ctx.PurgeQueue(w, r)
			break
		}
	case "ListDeadLetterSourceQueues":
		{
			ctx.ListDeadLetterSourceQueues(w, r)
			break
		}
	case "ChangeMessageVisibility":
		{
			ctx.ChangeMessageVisibility(w, r)
			break
		}
	default:
		{
			resp := toXMLErrorResponse("InvalidAction", fmt.Sprintf("The action %s is not valid for this endpoint.", r.FormValue("Action")), "")
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/xml")
			w.Write(resp)
			return
		}
	}
}
