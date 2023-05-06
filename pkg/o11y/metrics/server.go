package metrics

import (
	"context"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/charlie1404/vqs/pkg/o11y/logs"
)

type MetricsServer struct {
	httpServer *http.Server
}

func (s *MetricsServer) SetupInterruptListener() {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-stopCh
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logs.Logger.Warn().Msg("interrupt signal received, shutting down metrics server")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		logs.Logger.Fatal().Err(err).Msg("Unable to shutdown")
	}
}

func (s *MetricsServer) StartServer() {
	srvMux := http.NewServeMux()

	srvMux.HandleFunc("/debug/pprof/", pprof.Index)
	srvMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	srvMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	srvMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	srvMux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	srvMux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		Registry: registry,
	}))

	s.httpServer = &http.Server{
		Addr:    "127.0.0.1:1337",
		Handler: srvMux,
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

func New() *MetricsServer {
	initRegistry()
	s := &MetricsServer{}
	return s
}
