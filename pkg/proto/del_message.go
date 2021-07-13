package proto

import "net/http"

type DeleteMessagePayload struct {
	MessageId string
	QueueName string
}

func (ctx *Context) DeleteMessage(w http.ResponseWriter, req *http.Request) {

}
