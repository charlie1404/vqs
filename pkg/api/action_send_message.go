package api

import (
	"crypto/md5"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/charlie1404/vqs/pkg/app_errors"
	"github.com/charlie1404/vqs/pkg/o11y/logs"
	"github.com/charlie1404/vqs/pkg/storage"
	"github.com/charlie1404/vqs/pkg/utils"
)

type MessageAttribute struct {
	DataType string
	Value    []byte
}

type MessageAttributes map[string]MessageAttribute

type SendMessageInput struct {
	DelaySeconds            uint16            `validate:"min=0,max=900"`
	MessageBody             string            `validate:"required,min=1,max=262144"`
	QueueName               string            `validate:"required"`
	AccountId               string            `validate:"required"`
	MessageAttributes       MessageAttributes `validate:"required"`
	MessageSystemAttributes MessageAttributes ``
}

func parseSendMessageInput(form url.Values) (*SendMessageInput, error) {
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

	delaySeconds := utils.GetFormValueUint(form, "DelaySeconds", 0)

	messageBody := utils.GetFormValueString(form, "MessageBody")

	// Max 10 attributes are allowed
	messageAttributes := make(MessageAttributes)
	for i := 1; i <= 10; i++ {
		attribType := utils.GetFormValueString(form, fmt.Sprintf("MessageAttribute.%d.Value.DataType", i))
		if attribType != "String" && attribType != "Number" && attribType != "Binary" {
			continue
		}
		attribName := utils.GetFormValueString(form, fmt.Sprintf("MessageAttribute.%d.Name", i))
		attribValue := utils.GetFormValueString(form, fmt.Sprintf("MessageAttribute.%d.Value.%sValue", i, attribType))

		if len(attribName) < 1 ||
			len(attribName) > 256 {
			continue
		}

		messageAttributes[attribName] = MessageAttribute{
			DataType: attribType,
			Value:    []byte(attribValue),
		}
	}

	sendMessageInput := SendMessageInput{
		DelaySeconds:            uint16(delaySeconds),
		QueueName:               accountIdAndQueueName[1],
		MessageBody:             messageBody,
		AccountId:               accountIdAndQueueName[0],
		MessageAttributes:       messageAttributes,
		MessageSystemAttributes: make(MessageAttributes),
	}

	return &sendMessageInput, nil
}

func (appCtx *AppContext) SendMessage(w http.ResponseWriter, r *http.Request) {
	sendMessageInput, _ := parseSendMessageInput(r.Form)

	var queue *storage.Queue
	var err error

	queue, err = appCtx.queues.GetQueue(sendMessageInput.QueueName)
	if err != nil && err == app_errors.QueueNotExists {
		logs.Logger.Warn().Msg("Queue not preset, for push, creating one with defaults to recover")
		if queue, err = appCtx.queues.CreateQueue(sendMessageInput.QueueName, 0, 262144, 345600, 0, 30, &map[string]string{}); err == app_errors.CreateQueueQueueExists {
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

	message := storage.NewMessage(sendMessageInput.DelaySeconds, sendMessageInput.MessageBody)
	for k, v := range sendMessageInput.MessageAttributes {
		message.Attributes[k] = storage.Attribute{DataType: v.DataType, Value: v.Value}
	}
	for k, v := range sendMessageInput.MessageSystemAttributes {
		message.SystemAttributes[k] = storage.Attribute{DataType: v.DataType, Value: v.Value}
	}

	if err = queue.Push(message); err != nil {
		log.Println(err)
		resp := toXMLErrorResponse("UnknowError", "TODO !implement later.", "")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	resp := toSendMessageResponse(sendMessageInput.MessageBody)
	w.Write(resp)
}

type SendMessageResult struct {
	MessageId                    string
	MD5OfMessageBody             string
	MD5OfMessageAttributes       string
	MD5OfMessageSystemAttributes string
}

type SendMessageResponse struct {
	XMLName           xml.Name          `xml:"http://queue.amazonaws.com/doc/2012-11-05/ CreateQueueResponse"`
	SendMessageResult SendMessageResult `xml:"SendMessageResult"`
	ResponseMetadata  ResponseMetadata  `xml:"ResponseMetadata"`
}

func toSendMessageResponse(body string) []byte {
	resp := SendMessageResponse{
		SendMessageResult: SendMessageResult{
			MD5OfMessageBody:             fmt.Sprintf("%x", md5.Sum([]byte(body))),
			MD5OfMessageAttributes:       "3ae8f24a165a8cedc005670c81a27295",
			MD5OfMessageSystemAttributes: "3ae8f24a165a8cedc005670c81a27295",
			MessageId:                    utils.GenerateUUIDLikeId(),
		},
		ResponseMetadata: NewResponseMetadata(),
	}

	response, _ := xml.Marshal(resp)
	response = append([]byte("<?xml version=\"1.0\"?>"), response...)

	return response
}
