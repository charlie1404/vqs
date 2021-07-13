package proto

import (
	"encoding/json"
	"net/http"

	"github.com/charlie1404/vqueue/pkg/message"
	"github.com/charlie1404/vqueue/pkg/utils"
)

type SendMessagePayload struct {
	MessageBody       message.Message
	QueueName         string
	DelaySeconds      uint16
	MessageAttributes message.MessageAttributes
}

func (ctx *Context) SendMessage(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t SendMessagePayload
	err := decoder.Decode(&t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Malformed Request"))
		return
	}

	message_id := []byte(utils.GenRandomId())
	ctx.buffer.Insert(message_id)

	w.Write(message_id)
}
