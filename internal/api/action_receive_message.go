package api

import (
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/valyala/fasthttp"

	"github.com/charlie1404/vqs/internal/app_errors"
	"github.com/charlie1404/vqs/internal/o11y/logs"
	"github.com/charlie1404/vqs/internal/storage"
	"github.com/charlie1404/vqs/internal/utils"
)

type ReceiveMessageInput struct {
	AttributeName        string `validate:"required"`
	MaxNumberOfMessages  uint8  `validate:"min=1,max=10"`
	QueueName            string `validate:"required"`
	AccountId            string `validate:"required"`
	MessageAttributeName string `validate:"required"`
}

func parseReceiveMessageInput(form FormValues) (*ReceiveMessageInput, error) {
	parsedQueueUrl, err := url.Parse(form["QueueUrl"])
	if err != nil {
		return nil, err // return queue not found error
	}

	accountIdAndQueueName := strings.Split(parsedQueueUrl.Path[1:], "/")
	if len(accountIdAndQueueName) != 2 {
		return nil, errors.New("invalid queueUrl") // // return queue not found error
	}

	receiveMessageInput := ReceiveMessageInput{
		QueueName: accountIdAndQueueName[1],
		AccountId: accountIdAndQueueName[0],
	}

	return &receiveMessageInput, nil
}

func (appCtx *AppContext) ReceiveMessage(ctx *fasthttp.RequestCtx) {
	receiveMessageInput, _ := parseReceiveMessageInput(ctx.UserValue("body").(FormValues))

	x := StreamWriter{ctx}

	var queue *storage.Queue
	var err error

	queue, err = appCtx.queues.GetQueue(receiveMessageInput.QueueName)
	if err != nil && err == app_errors.QueueNotExists {
		logs.Logger.Warn().Msg("QueueNotPresentFault, creating one with defaults to recover")
		if queue, err = appCtx.queues.CreateQueue(receiveMessageInput.QueueName, 0, 262144, 345600, 0, 30, &[][2]string{}); err == app_errors.CreateQueueQueueExists {
			logs.Logger.Info().Msg("Queue created by another proc")
			err = nil
		}
	}

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		appCtx.templates.ExecuteTemplate(x, "error.tpl", NewVQSError("UnknowError", "TODO !implement later.", ""))
		return
	}

	msg, err := queue.Pop()
	if err != nil {
		log.Println(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		appCtx.templates.ExecuteTemplate(x, "error.tpl", NewVQSError("UnknowError", "TODO !implement later.", ""))
		return
	}

	if msg == nil { // TODO: remove me
		msg = &storage.Message{
			Id:   "todo1234567890asdfghjkl",
			Body: "todo1234567890asdfghjkl",
		}
	}

	val := &ReceiveMessage{
		MessageId:              msg.Id,
		Body:                   msg.Body,
		ReceiptHandle:          "todo234567890sdertyuiop",
		MD5OfBody:              fmt.Sprintf("%x", md5.Sum([]byte(msg.Body))),
		MD5OfMessageAttributes: fmt.Sprintf("%x", md5.Sum([]byte(""))),
		RequestId:              utils.GenerateUUIDLikeId(),
	}

	appCtx.templates.ExecuteTemplate(x, "receive_message.tpl", val)
}

type ReceiveMessage struct {
	MessageId              string
	ReceiptHandle          string
	Body                   string
	MD5OfBody              string
	MD5OfMessageAttributes string
	RequestId              string
	// Attributes   []Attribute
	// MessageAttributes []MessageAttribute
}
