package api

import (
	"log"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/charlie1404/vqueue/pkg/storage"
)

type ApiApp struct {
	httpServer *http.Server
}

func New() *ApiApp {
	s := &ApiApp{}
	return s
}

// func (s *ApiApp) Shutdown() {
// 	if s.httpServer != nil {
// 		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
// 		err := s.httpServer.Shutdown(ctx)
// 		if err != nil {
// 			log.Println("Failed to shutdown http server gracefully")
// 			log.Fatalln(err)
// 		} else {
// 			s.httpServer = nil
// 		}
// 	}
// }

func (s *ApiApp) StartServer() {
	validator := newValidator()
	validator.registerCustomValidations()

	appCtx := &AppContext{
		queues:    storage.NewQueues(),
		validator: validator,
	}

	srvMux := http.NewServeMux()

	srvMux.HandleFunc("/", timeoutMiddleware(appCtx.requestHandler))

	srvMux.HandleFunc("/debug/pprof/", pprof.Index)
	srvMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	srvMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	srvMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	srvMux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	s.httpServer = &http.Server{
		Addr:              "localhost:3344",
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           srvMux,

		// BaseContext: func(l net.Listener) context.Context {
		// 	ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
		// 	return ctx
		// },
	}

	if err := s.httpServer.ListenAndServe(); err != nil {
		log.Println("Http Server stopped unexpected")
		log.Fatalln(err)
	}
}
