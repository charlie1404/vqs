package proto

import (
	"encoding/json"
	"net/http"
)

type ReceiveMessagePayload struct {
	QueueName string
	// MaxNumberOfMessages uint8
	// WaitTimeSeconds     uint16
}

func (ctx *Context) ReceiveMessage(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t ReceiveMessagePayload
	err := decoder.Decode(&t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Malformed Request"))
		return
	}

	message_id := ctx.buffer.Pop()

	w.Write(message_id)
}
