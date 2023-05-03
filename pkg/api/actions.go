package api

import (
	"net/http"
)

type MessageAttribute struct {
	DataType string
	Value    []byte
}

type MessageAttributes map[string]MessageAttribute

func (appCtx *AppContext) DeleteMessage(w http.ResponseWriter, r *http.Request) {
}

func (appCtx *AppContext) ListQueues(w http.ResponseWriter, r *http.Request) {
}

func (appCtx *AppContext) DeleteMessageBatch(w http.ResponseWriter, r *http.Request) {
}

func (appCtx *AppContext) SendMessageBatch(w http.ResponseWriter, r *http.Request) {
}

func (appCtx *AppContext) DeleteQueue(w http.ResponseWriter, r *http.Request) {
}

func (appCtx *AppContext) TagQueue(w http.ResponseWriter, r *http.Request) {
}

func (appCtx *AppContext) ListQueueTags(w http.ResponseWriter, r *http.Request) {
}

func (appCtx *AppContext) UntagQueue(w http.ResponseWriter, r *http.Request) {
}

func (appCtx *AppContext) SetQueueAttributes(w http.ResponseWriter, r *http.Request) {
}

func (appCtx *AppContext) GetQueueAttributes(w http.ResponseWriter, r *http.Request) {
}

func (appCtx *AppContext) GetQueueUrl(w http.ResponseWriter, r *http.Request) {
}

func (appCtx *AppContext) PurgeQueue(w http.ResponseWriter, r *http.Request) {
}

func (appCtx *AppContext) ListDeadLetterSourceQueues(w http.ResponseWriter, r *http.Request) {
}

func (appCtx *AppContext) ChangeMessageVisibility(w http.ResponseWriter, r *http.Request) {
}
