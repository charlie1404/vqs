package proto

import (
	"net/http"
	"net/http/pprof"

	"github.com/charlie1404/vqueue/pkg/buffer"
)

type Context struct {
	buffer *buffer.Buffer
}

func StartServe(buffer *buffer.Buffer) {
	ctx := &Context{
		buffer: buffer,
	}

	router := http.NewServeMux()
	router.HandleFunc("/SendMessage", ctx.SendMessage)
	router.HandleFunc("/ReceiveMessage", ctx.ReceiveMessage)
	router.HandleFunc("/DeleteMessage", ctx.ReceiveMessage)

	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)

	err := http.ListenAndServe(":3344", router)
	if err != nil {
		panic(err)
	}
}
