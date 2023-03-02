package storage

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"

	"github.com/charlie1404/vqs/pkg/utils"
)

type Attribute struct {
	DataType string
	Value    []byte
}

type Message struct {
	Id               string
	Body             string
	Attributes       map[string]Attribute
	SystemAttributes map[string]Attribute
	DelaySeconds     uint16
}

func NewMessage(delaySeconds uint16, body string) *Message {
	return &Message{
		Id:               utils.GenerateUUIDLikeId(),
		DelaySeconds:     delaySeconds,
		Body:             body,
		Attributes:       make(map[string]Attribute),
		SystemAttributes: make(map[string]Attribute),
	}
}

func serializeMessage(m interface{}) ([]byte, error) {
	// for debugging purposes use json encoding instead of gob+gzip

	var gobBuffer bytes.Buffer
	var gzipBuffer bytes.Buffer

	enc := gob.NewEncoder(&gobBuffer)
	enc.Encode(m)

	compressor := gzip.NewWriter(&gzipBuffer)
	_, err := compressor.Write(gobBuffer.Bytes())
	if err != nil {
		return nil, err
	}
	compressor.Close()

	return gzipBuffer.Bytes(), nil
}

func deserializeMessage(b []byte) (*Message, error) {
	// for debugging purposes use json encoding instead of gob+gzip

	var gobBuffer bytes.Buffer
	var gzipBuffer bytes.Buffer

	gzipBuffer.Write(b)
	decompressor, err := gzip.NewReader(&gzipBuffer)
	if err != nil {
		return nil, err
	}
	_, err = gobBuffer.ReadFrom(decompressor)
	if err != nil {
		return nil, err
	}
	decompressor.Close()

	dec := gob.NewDecoder(&gobBuffer)
	var m Message
	err = dec.Decode(&m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}
