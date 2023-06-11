package api

import (
	"encoding/xml"

	"github.com/charlie1404/vqs/internal/utils"
)

type VQSError struct {
	Code      string
	Detail    string
	Message   string
	Type      string
	RequestId string
}

func NewVQSError(code, message, detail string) *VQSError {
	return &VQSError{
		Code:      code,
		Message:   message,
		Detail:    detail,
		Type:      "Sender",
		RequestId: utils.GenerateUUIDLikeId(),
	}
}

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

func NewResponseMetadata() ResponseMetadata {
	return ResponseMetadata{
		RequestId: utils.GenerateUUIDLikeId(),
	}
}

func toXMLErrorResponse(code, message, detail string) []byte {
	errResp := ErrorResponse{
		Error: ErrorResponseError{
			Code:    code,
			Message: message,
			Detail:  "",
			Type:    "Sender",
		},
		RequestId: utils.GenerateUUIDLikeId(),
	}

	response, _ := xml.Marshal(errResp)
	response = append([]byte("<?xml version=\"1.0\"?>"), response...)

	return response
}
