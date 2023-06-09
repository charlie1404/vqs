package api

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/charlie1404/vqs/pkg/app_errors"
	"github.com/charlie1404/vqs/pkg/o11y/logs"
	"github.com/charlie1404/vqs/pkg/storage"
	"github.com/charlie1404/vqs/pkg/utils"
	"github.com/valyala/fasthttp"
)

type SendMessageInput struct {
	DelaySeconds            uint16            `validate:"min=0,max=900"`
	MessageBody             string            `validate:"required,min=1,max=262144"`
	QueueName               string            `validate:"required"`
	AccountId               string            `validate:"required"`
	MessageAttributes       MessageAttributes `validate:"required"`
	MessageSystemAttributes MessageAttributes ``
}

func parseSendMessageInput(form FormValues) (*SendMessageInput, error) {
	parsedQueueUrl, err := url.Parse(form["QueueUrl"])
	if err != nil {
		return nil, err // return queue not found error
	}

	accountIdAndQueueName := strings.Split(parsedQueueUrl.Path[1:], "/")
	if len(accountIdAndQueueName) != 2 {
		return nil, errors.New("invalid queueUrl") // return queue not found error
	}

	delaySeconds, err := strconv.ParseUint(form["DelaySeconds"], 10, 16)
	if err != nil {
		delaySeconds = 0
	}

	// Max 10 attributes are allowed
	messageAttributes := make(MessageAttributes)
	for i := 1; i <= 10; i++ {
		attribType := form[fmt.Sprintf("MessageAttribute.%d.Value.DataType", i)]
		if attribType != "String" && attribType != "Number" && attribType != "Binary" {
			continue
		}

		var attribValue []byte

		attribName := form[fmt.Sprintf("MessageAttribute.%d.Name", i)]
		attribStringValue := form[fmt.Sprintf("MessageAttribute.%d.Value.StringValue", i)]
		attribBinaryValue := form[fmt.Sprintf("MessageAttribute.%d.Value.BinaryValue", i)]

		// string value has higher priority
		if attribStringValue != "" {
			attribValue = []byte(attribStringValue) // TODO: check if data type is binary, then base64 decode is required
		}

		if attribStringValue == "" && attribBinaryValue != "" {
			data, err := base64.StdEncoding.DecodeString(attribBinaryValue)
			if err == nil { // ignore if base64 decode fails
				attribValue = data
			}
		}

		if len(attribName) < 1 ||
			len(attribName) > 256 ||
			len(attribValue) < 1 {
			logs.Logger.Warn().Msg("Invalid attribute")
			continue
		}

		messageAttributes[attribName] = MessageAttribute{
			DataType: attribType,
			Value:    attribValue,
		}
	}

	sendMessageInput := SendMessageInput{
		DelaySeconds:            uint16(delaySeconds),
		QueueName:               accountIdAndQueueName[1],
		MessageBody:             form["MessageBody"],
		AccountId:               accountIdAndQueueName[0],
		MessageAttributes:       messageAttributes,
		MessageSystemAttributes: make(MessageAttributes),
	}

	return &sendMessageInput, nil
}

func (appCtx *AppContext) SendMessage(ctx *fasthttp.RequestCtx) {
	sendMessageInput, _ := parseSendMessageInput(ctx.UserValue("body").(FormValues))

	var queue *storage.Queue
	var err error

	queue, err = appCtx.queues.GetQueue(sendMessageInput.QueueName)
	if err != nil && err == app_errors.QueueNotExists {
		logs.Logger.Warn().Msg("QueueNotPresentFault, creating one with defaults to recover")

		if queue, err = appCtx.queues.CreateDefaultQueue(sendMessageInput.QueueName); err == app_errors.CreateQueueQueueExists {
			logs.Logger.Info().Msg("Queue created by another proc")
			err = nil
		}
	}

	if err != nil {
		log.Println(err)
		resp := toXMLErrorResponse("UnknowError", "TODO !implement later.", "")
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody(resp)
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
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody(resp)
		return
	}

	resp := toSendMessageResponse(sendMessageInput.MessageBody)
	ctx.SetBody(resp)
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
