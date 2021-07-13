package api

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/xml"
	"fmt"
)

type ErrorResponseError struct {
	Code    string `xml:"Code"`
	Detail  string `xml:"Detail"`
	Message string `xml:"Message"`
	Type    string `xml:"Type"`
}

type ErrorResponse struct {
	XMLName   xml.Name           `xml:"http://queue.amazonaws.com/doc/2012-11-05/ ErrorResponse"`
	Error     ErrorResponseError `xml:"Error"`
	RequestId string             `xml:"RequestId"`
}

type ResponseMetadata struct {
	RequestId string `xml:"RequestId"`
}

func debugPrint(resp interface{}) {
	response, _ := xml.MarshalIndent(resp, "", "  ")
	response = append([]byte("<?xml version=\"1.0\"?>\n"), response...)
	fmt.Println(string(response))
}

func generateRequestId() string {
	b := make([]byte, 10)
	rand.Read(b)
	randomHash := md5.Sum(b)

	return fmt.Sprintf("%x-%x-%x-%x-%x", randomHash[0:4], randomHash[4:6], randomHash[6:8], randomHash[8:10], randomHash[10:])
}

func toXMLErrorResponse(code, message, detail string) []byte {
	errResp := ErrorResponse{
		Error: ErrorResponseError{
			Code:    code,
			Message: message,
			Detail:  "",
			Type:    "Sender",
		},
		RequestId: generateRequestId(),
	}

	response, _ := xml.Marshal(errResp)
	response = append([]byte("<?xml version=\"1.0\"?>"), response...)

	return response
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
		ResponseMetadata: ResponseMetadata{
			RequestId: generateRequestId(),
		},
	}

	response, _ := xml.Marshal(resp)
	response = append([]byte("<?xml version=\"1.0\"?>"), response...)

	return response
}
