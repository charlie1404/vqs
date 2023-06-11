package api

import (
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"

	"github.com/charlie1404/vqs/internal/o11y/logs"
	"github.com/charlie1404/vqs/internal/storage"
	"github.com/valyala/fasthttp"
)

type AppContext struct {
	queues    *storage.Queues
	templates *template.Template
}

type ApiApp struct {
	httpServer *fasthttp.Server
	appCtx     *AppContext
}

func (s *ApiApp) SetupInterruptListener() {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-stopCh
	logs.Logger.Info().Msg("interrupt signal received, shutting down metrics server")

	s.httpServer.DisableKeepalive = true

	time.Sleep(1 * time.Second)

	if err := s.httpServer.Shutdown(); err != nil {
		logs.Logger.Fatal().Err(err).Msg("Unable to shutdown")
	}
}

func (s *ApiApp) StartServer() {
	templates := template.Must(template.ParseGlob("internal/templates/*"))
	middleware := Middleware{templates: templates}

	s.appCtx = &AppContext{
		queues:    storage.LoadQueues(),
		templates: templates,
	}

	s.httpServer = &fasthttp.Server{
		Handler:              middleware.WrapHandler(s.appCtx.requestHandler),
		ReadTimeout:          5 * time.Second,
		WriteTimeout:         5 * time.Second,
		IdleTimeout:          30 * time.Second,
		MaxConnsPerIP:        500,
		MaxRequestsPerConn:   500,
		MaxKeepaliveDuration: 5 * time.Second,
	}

	go func() {
		logs.Logger.Info().Msg("Starting http server")
		err := s.httpServer.ListenAndServe("127.0.0.1:3344")
		if err != nil {
			logs.Logger.Fatal().Err(err).Msg("Http Server stopped unexpected")
		}
	}()
}

func (s *ApiApp) CloseQueues() {
	s.appCtx.queues.CloseQueues()
}

func New() *ApiApp {
	s := &ApiApp{}
	return s
}
