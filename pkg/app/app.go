package app

import (
	"github.com/charlie1404/vqueue/pkg/buffer"
	"github.com/charlie1404/vqueue/pkg/proto"
)

func New() {
	quequBuffer, err := buffer.New("sgtest", 64)
	if err != nil {
		panic(err)
	}

	proto.StartServe(quequBuffer)
}
