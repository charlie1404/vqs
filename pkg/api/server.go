package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charlie1404/vqs/pkg/o11y/logs"
	"github.com/charlie1404/vqs/pkg/storage"
)

type AppContext struct {
	queues    *storage.Queues
	validator ApiValidator
}

type ApiApp struct {
	httpServer *http.Server
	appCtx     *AppContext
}

func (s *ApiApp) SetupInterruptListener() {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-stopCh
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	logs.Logger.Warn().Msg("interrupt signal received, shutting down metrics server")

	// time.Sleep(10 * time.Second)

	if err := s.httpServer.Shutdown(ctx); err != nil {
		logs.Logger.Fatal().Err(err).Msg("Unable to shutdown")
	}
}

func (s *ApiApp) StartServer() {
	validator := newValidator()
	validator.registerCustomValidations()

	s.appCtx = &AppContext{
		queues:    storage.LoadQueues(),
		validator: validator,
	}

	srvMux := http.NewServeMux()

	srvMux.Handle("/", logRequestHandler(
		timeoutMiddleware(
			http.HandlerFunc(
				s.appCtx.requestHandler,
			),
		),
	),
	)

	s.httpServer = &http.Server{
		Addr:              "127.0.0.1:3344",
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

	go func() {
		logs.Logger.Info().Msg("Starting http server")
		err := s.httpServer.ListenAndServe()
		if err == http.ErrServerClosed {
			logs.Logger.Warn().Msg("Http Server stopped")
			return
		}
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
