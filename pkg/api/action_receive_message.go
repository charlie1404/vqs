package api

import (
	"encoding/xml"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/charlie1404/vqs/pkg/app_errors"
	"github.com/charlie1404/vqs/pkg/o11y/logs"
	"github.com/charlie1404/vqs/pkg/storage"
	"github.com/charlie1404/vqs/pkg/utils"
)

type ReceiveMessageInput struct {
	AttributeName        string `validate:"required"`
	MaxNumberOfMessages  uint8  `validate:"min=1,max=10"`
	QueueName            string `validate:"required"`
	AccountId            string `validate:"required"`
	MessageAttributeName string `validate:"required"`
}

func parseReceiveMessageInput(form url.Values) (*ReceiveMessageInput, error) {
	queueUrl := utils.GetFormValueString(form, "QueueUrl")

	parsedQueueUrl, err := url.Parse(queueUrl)
	if err != nil {
		// todo retrun standard error
		return nil, err
	}

	accountIdAndQueueName := strings.Split(parsedQueueUrl.Path[1:], "/")
	if len(accountIdAndQueueName) != 2 {
		// todo retrun standard error
		return nil, errors.New("invalid queueUrl")
	}

	receiveMessageInput := ReceiveMessageInput{
		QueueName: accountIdAndQueueName[1],
		AccountId: accountIdAndQueueName[0],
	}

	return &receiveMessageInput, nil
}

func (appCtx *AppContext) ReceiveMessage(w http.ResponseWriter, r *http.Request) {
	receiveMessageInput, _ := parseReceiveMessageInput(r.Form)
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
		log.Println(err)
		resp := toXMLErrorResponse("UnknowError", "TODO !implement later.", "")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	msg, err := queue.Pop()
	if err != nil {
		resp := toXMLErrorResponse("UnknowError", "TODO !implement later.", "")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	resp := toReceiveMessageResponse(msg)
	w.Write(resp)
}

type Message struct {
	MessageId              string
	ReceiptHandle          string
	Body                   string
	MD5OfBody              string
	MD5OfMessageAttributes string
	// Attributes   []Attribute
	// MessageAttributes []MessageAttribute
}

type ReceiveMessageResult struct {
	Message []Message
}

type ReceiveMessageResponse struct {
	XMLName              xml.Name             `xml:"http://queue.amazonaws.com/doc/2012-11-05/ CreateQueueResponse"`
	ReceiveMessageResult ReceiveMessageResult `xml:"ReceiveMessageResult"`
	ResponseMetadata     ResponseMetadata     `xml:"ResponseMetadata"`
}

func toReceiveMessageResponse(msg *storage.Message) []byte {
	resp := ReceiveMessageResponse{
		ReceiveMessageResult: ReceiveMessageResult{
			Message: []Message{}, // todo
		},
		ResponseMetadata: NewResponseMetadata(),
	}

	response, _ := xml.Marshal(resp)
	response = append([]byte("<?xml version=\"1.0\"?>"), response...)

	return response
}
