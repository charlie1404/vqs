package api

import (
	"fmt"
	"net/http"
)

func (ctx *AppContext) requestHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	w.Header().Set("Content-Type", "application/xml")

	if r.Method != "POST" {
		resp := toXMLErrorResponse("UnsupportedMethod", fmt.Sprintf("(%s) method is not supported", r.Method), "")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	if r.FormValue("Version") != "2012-11-05" {
		resp := toXMLErrorResponse("NoSuchVersion", fmt.Sprintf("The requested version ( %s ) is not valid.", r.FormValue("Version")), "")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	switch r.FormValue("Action") {
	case "ListQueues":
		ctx.ListQueues(w, r)
		break
	case "CreateQueue":
		ctx.CreateQueue(w, r)
		break
	case "SendMessage":
		ctx.SendMessage(w, r)
		break
	case "ReceiveMessage":
		ctx.ReceiveMessage(w, r)
		break
	case "DeleteMessage":
		ctx.DeleteMessage(w, r)
		break
	case "DeleteMessageBatch":
		ctx.DeleteMessageBatch(w, r)
		break
	case "SendMessageBatch":
		ctx.SendMessageBatch(w, r)
		break
	case "DeleteQueue":
		ctx.DeleteQueue(w, r)
		break
	case "TagQueue":
		ctx.TagQueue(w, r)
		break
	case "ListQueueTags":
		ctx.ListQueueTags(w, r)
		break
	case "UntagQueue":
		ctx.UntagQueue(w, r)
		break
	case "SetQueueAttributes":
		ctx.SetQueueAttributes(w, r)
		break
	case "GetQueueAttributes":
		ctx.GetQueueAttributes(w, r)
		break
	case "GetQueueUrl":
		ctx.GetQueueUrl(w, r)
		break
	case "PurgeQueue":
		ctx.PurgeQueue(w, r)
		break
	case "ListDeadLetterSourceQueues":
		ctx.ListDeadLetterSourceQueues(w, r)
		break
	case "ChangeMessageVisibility":
		ctx.ChangeMessageVisibility(w, r)
		break
	default:
		resp := toXMLErrorResponse("InvalidAction", fmt.Sprintf("The action %s is not valid for this endpoint.", r.FormValue("Action")), "")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}
}
