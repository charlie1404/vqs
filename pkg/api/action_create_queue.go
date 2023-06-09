package api

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/valyala/fasthttp"

	"github.com/charlie1404/vqs/pkg/o11y/logs"
	"github.com/charlie1404/vqs/pkg/storage"
)

type QueueAttributes struct {
	DelaySeconds                  uint16 `validate:"min=0,max=900"`
	MaximumMessageSize            uint32 `validate:"min=1024,max=262144"`
	MessageRetentionPeriod        uint32 `validate:"min=60,max=1209600"`
	ReceiveMessageWaitTimeSeconds uint16 `validate:"min=0,max=20"`
	VisibilityTimeout             uint16 `validate:"min=0,max=43200"`
}

type CreateQueueInput struct {
	QueueName  string          `validate:"required,min=5,max=50"`
	Attributes QueueAttributes `validate:"required"`
	Tags       *[][2]string    ``
}

func parseCreateQueueInput(form FormValues) *CreateQueueInput {
	var (
		delaySeconds                  uint16 = 0
		maximumMessageSize            uint32 = 262144
		messageRetentionPeriod        uint32 = 345600
		receiveMessageWaitTimeSeconds uint16 = 0
		visibilityTimeout             uint16 = 30
	)

	// Attributes
	for i := 1; i <= 8; i++ {
		attribName := form[fmt.Sprintf("Attribute.%d.Name", i)]
		val := form[fmt.Sprintf("Attribute.%d.Value", i)]

		switch attribName {
		case "DelaySeconds":
			if intVal, err := strconv.ParseUint(val, 10, 16); err == nil {
				delaySeconds = uint16(intVal)
			}
		case "MaximumMessageSize":
			if intVal, err := strconv.ParseUint(val, 10, 32); err == nil {
				maximumMessageSize = uint32(intVal)
			}
		case "MessageRetentionPeriod":
			if intVal, err := strconv.ParseUint(val, 10, 32); err == nil {
				messageRetentionPeriod = uint32(intVal)
			}
		case "ReceiveMessageWaitTimeSeconds":
			if intVal, err := strconv.ParseUint(val, 10, 16); err == nil {
				receiveMessageWaitTimeSeconds = uint16(intVal)
			}
		case "VisibilityTimeout":
			if intVal, err := strconv.ParseUint(val, 10, 16); err == nil {
				visibilityTimeout = uint16(intVal)
			}
		}
	}

	// Tags Limitations
	// - Tags after 50 will be discarded
	// - Max size of tag storage is meta file size - meta header
	// - MaxKeyLength is 128 in UTF-8. The tag Key must not be empty or null.
	// - MaximumTagValueLength is 256 in UTF-8. The tag Value may be empty or null.
	tags := [][2]string{}
	tagsSize := 0
	for i := 1; i <= 50; i++ {
		tagName := form[fmt.Sprintf("Tag.%d.Key", i)]
		tagValue := form[fmt.Sprintf("Tag.%d.Value", i)]

		if len(tagName) == 0 ||
			len(tagValue) == 0 ||
			utf8.RuneCountInString(tagName) > 128 ||
			utf8.RuneCountInString(tagValue) > 256 {
			continue
		}

		maxSize := int(storage.META_FILE_SIZE - storage.META_FILE_META_DATA_SIZE)
		if tagsSize += len(tagName) + len(tagValue) + 4; tagsSize > maxSize {
			break
		}

		tags = append(tags, [2]string{tagName, tagValue})
	}

	createQueueInput := CreateQueueInput{
		QueueName: form["QueueName"],
		Attributes: QueueAttributes{
			DelaySeconds:                  delaySeconds,
			MaximumMessageSize:            maximumMessageSize,
			MessageRetentionPeriod:        messageRetentionPeriod,
			ReceiveMessageWaitTimeSeconds: receiveMessageWaitTimeSeconds,
			VisibilityTimeout:             visibilityTimeout,
		},
		Tags: &tags,
	}

	return &createQueueInput
}

func (appCtx *AppContext) CreateQueue(ctx *fasthttp.RequestCtx) {
	ip := parseCreateQueueInput(ctx.UserValue("body").(FormValues))
	fmt.Printf("%+v\n", ip)

	// TODO: Validate input
	// if err := appCtx.validator.validateCreateQueueInput(ip); err != nil {
	// 	resp := toXMLErrorResponse("InvalidAttributeValue", "Invalid value for some parameter.", "")
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Header().Set("Content-Type", "application/xml")
	// 	w.Write(resp)
	// 	return
	// }

	if _, err := appCtx.queues.CreateQueue(ip.QueueName, ip.Attributes.DelaySeconds, ip.Attributes.MaximumMessageSize, ip.Attributes.MessageRetentionPeriod, ip.Attributes.ReceiveMessageWaitTimeSeconds, ip.Attributes.VisibilityTimeout, ip.Tags); err != nil {
		logs.Logger.Error().Err(err).Msg("CreateQueue")
		resp := toXMLErrorResponse("InternalServerError", "Todo return better errors", "")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody(resp)
		return
	}

	resp := toCreateQueueResponse(ip.QueueName)
	ctx.SetBody(resp)
}

type CreateQueueResult struct {
	QueueUrl string `xml:"QueueUrl"`
}

type CreateQueueResponse struct {
	XMLName           xml.Name          `xml:"http://queue.amazonaws.com/doc/2012-11-05/ CreateQueueResponse"`
	CreateQueueResult CreateQueueResult `xml:"CreateQueueResult"`
	ResponseMetadata  ResponseMetadata  `xml:"ResponseMetadata"`
}

func toCreateQueueResponse(queueName string) []byte {
	resp := CreateQueueResponse{
		CreateQueueResult: CreateQueueResult{
			QueueUrl: queueName,
		},
		ResponseMetadata: NewResponseMetadata(),
	}

	response, _ := xml.Marshal(resp)
	response = append([]byte("<?xml version=\"1.0\"?>"), response...)

	return response
}
