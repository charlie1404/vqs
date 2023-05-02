package api

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"unicode/utf8"

	"github.com/charlie1404/vqs/pkg/o11y/logs"
	"github.com/charlie1404/vqs/pkg/storage"
	"github.com/charlie1404/vqs/pkg/utils"
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

func parseCreateQueueInput(form url.Values) *CreateQueueInput {
	queueName := utils.GetFormValueString(form, "QueueName")
	var delaySeconds uint = 0
	var maximumMessageSize uint = 262144
	var messageRetentionPeriod uint = 345600
	var receiveMessageWaitTimeSeconds uint = 0
	var visibilityTimeout uint = 30

	// Attributes
	for i := 1; i <= 8; i++ {
		attribName := utils.GetFormValueString(form, fmt.Sprintf("Attribute.%d.Name", i))
		switch attribName {
		case "DelaySeconds":
			delaySeconds = utils.GetFormValueUint(form, fmt.Sprintf("Attribute.%d.Value", i), delaySeconds)
			continue
		case "MaximumMessageSize":
			maximumMessageSize = utils.GetFormValueUint(form, fmt.Sprintf("Attribute.%d.Value", i), maximumMessageSize)
			continue
		case "MessageRetentionPeriod":
			messageRetentionPeriod = utils.GetFormValueUint(form, fmt.Sprintf("Attribute.%d.Value", i), messageRetentionPeriod)
			continue
		case "ReceiveMessageWaitTimeSeconds":
			receiveMessageWaitTimeSeconds = utils.GetFormValueUint(form, fmt.Sprintf("Attribute.%d.Value", i), receiveMessageWaitTimeSeconds)
			continue
		case "VisibilityTimeout":
			visibilityTimeout = utils.GetFormValueUint(form, fmt.Sprintf("Attribute.%d.Value", i), visibilityTimeout)
			continue
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
		tagName := utils.GetFormValueString(form, fmt.Sprintf("Tag.%d.Key", i))
		tagValue := utils.GetFormValueString(form, fmt.Sprintf("Tag.%d.Value", i))

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
		QueueName: queueName,
		Attributes: QueueAttributes{
			DelaySeconds:                  uint16(delaySeconds),
			MaximumMessageSize:            uint32(maximumMessageSize),
			MessageRetentionPeriod:        uint32(messageRetentionPeriod),
			ReceiveMessageWaitTimeSeconds: uint16(receiveMessageWaitTimeSeconds),
			VisibilityTimeout:             uint16(visibilityTimeout),
		},
		Tags: &tags,
	}

	return &createQueueInput
}

func (appCtx *AppContext) CreateQueue(w http.ResponseWriter, r *http.Request) {
	ip := parseCreateQueueInput(r.Form)

	if err := appCtx.validator.validateCreateQueueInput(ip); err != nil {
		resp := toXMLErrorResponse("InvalidAttributeValue", "Invalid value for some parameter.", "")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/xml")
		w.Write(resp)
		return
	}

	if _, err := appCtx.queues.CreateQueue(ip.QueueName, ip.Attributes.DelaySeconds, ip.Attributes.MaximumMessageSize, ip.Attributes.MessageRetentionPeriod, ip.Attributes.ReceiveMessageWaitTimeSeconds, ip.Attributes.VisibilityTimeout, ip.Tags); err != nil {
		logs.Logger.Error().Err(err).Msg("CreateQueue")
		resp := toXMLErrorResponse("InternalServerError", "Todo return better errors", "")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	resp := toCreateQueueResponse(ip.QueueName)
	w.Write(resp)
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
