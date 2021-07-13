package api

import (
	"fmt"
	"net/url"
	"strconv"
	"unicode/utf8"

	"github.com/charlie1404/vqueue/pkg/constants"
	"github.com/charlie1404/vqueue/pkg/storage"
)

func getFormValueString(form url.Values, key string) string {
	if vals := form[key]; len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func getFormValueUint(form url.Values, key string, defaultVal uint) uint {
	val := ""
	if vals := form[key]; len(vals) > 0 {
		val = vals[0]
	}

	if intVal, err := strconv.Atoi(val); err != nil {
		return defaultVal
	} else {
		return uint(intVal)
	}
}

func parseCreateQueueInput(form url.Values) *storage.CreateQueueInput {
	queueName := getFormValueString(form, "QueueName")
	var delaySeconds uint = 0
	var maximumMessageSize uint = 262144
	var messageRetentionPeriod uint = 345600
	var receiveMessageWaitTimeSeconds uint = 0
	var visibilityTimeout uint = 30

	// Attributes
	for i := 1; i <= 8; i++ {
		attribName := getFormValueString(form, fmt.Sprintf("Attribute.%d.Name", i))
		switch attribName {
		case "DelaySeconds":
			delaySeconds = getFormValueUint(form, fmt.Sprintf("Attribute.%d.Value", i), delaySeconds)
			continue
		case "MaximumMessageSize":
			maximumMessageSize = getFormValueUint(form, fmt.Sprintf("Attribute.%d.Value", i), maximumMessageSize)
			continue
		case "MessageRetentionPeriod":
			messageRetentionPeriod = getFormValueUint(form, fmt.Sprintf("Attribute.%d.Value", i), messageRetentionPeriod)
			continue
		case "ReceiveMessageWaitTimeSeconds":
			receiveMessageWaitTimeSeconds = getFormValueUint(form, fmt.Sprintf("Attribute.%d.Value", i), receiveMessageWaitTimeSeconds)
			continue
		case "VisibilityTimeout":
			visibilityTimeout = getFormValueUint(form, fmt.Sprintf("Attribute.%d.Value", i), visibilityTimeout)
			continue
		}
	}

	createQueueInput := storage.CreateQueueInput{
		QueueName: queueName,
		Attributes: storage.QueueAttributes{
			DelaySeconds:                  delaySeconds,
			MaximumMessageSize:            maximumMessageSize,
			MessageRetentionPeriod:        messageRetentionPeriod,
			ReceiveMessageWaitTimeSeconds: receiveMessageWaitTimeSeconds,
			VisibilityTimeout:             visibilityTimeout,
		},
		Tags: make(storage.QueueTags),
	}

	// Tags Limitations
	// - Tags after 50 will be discarded
	// - Max size of tag storage 16350
	// - MaxKeyLength is 128 in UTF-8. The tag Key must not be empty or null.
	// - MaximumTagValueLength is 256 in UTF-8. The tag Value may be empty or null.

	tagsSize := 0
	for i := 1; i <= 50; i++ {
		tagName := getFormValueString(form, fmt.Sprintf("Tag.%d.Key", i))
		tagValue := getFormValueString(form, fmt.Sprintf("Tag.%d.Value", i))

		if len(tagName) == 0 ||
			len(tagValue) == 0 ||
			utf8.RuneCountInString(tagName) > 128 ||
			utf8.RuneCountInString(tagValue) > 256 {
			continue
		}

		if tagsSize += len(tagName) + len(tagValue) + 2; tagsSize > constants.MAX_HEADER_TAGS_SIZE {
			break
		}

		createQueueInput.Tags[tagName] = tagValue
	}

	return &createQueueInput
}
