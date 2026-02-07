package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/artarts36/go-metrics"

	"github.com/artarts36/ip-scanner/internal/port/grpc/app"
)

const initAppTimeout = 5 * time.Minute

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	ctx, cancel := context.WithCancel(context.Background())

	slog.InfoContext(ctx, "[main] read config")

	cfg, err := app.ReadConfig("IPSCANNER_")
	if err != nil {
		slog.ErrorContext(ctx, "[main] failed to read config", slog.Any("err", err))
		os.Exit(1)
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.Log.Level,
	})))

	metricsServer := metrics.NewServer(cfg.Metrics.ToGoMetrics())

	initAppCtx, cancelAppCtx := context.WithTimeout(ctx, initAppTimeout)

	application, err := app.NewApp(initAppCtx, cfg)
	if err != nil {
		cancelAppCtx()
		slog.ErrorContext(ctx, "[main] failed to read config", slog.Any("err", err))
		os.Exit(1)
	}

	defer cancelAppCtx()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.Info("serve application")

		if runErr := application.Run(ctx); runErr != nil {
			slog.
				With(slog.Any("err", runErr)).
				ErrorContext(ctx, "failed to serve application")
		}
	}()

	wg.Add(1)
	go func(s *metrics.Server) {
		slog.Info("serve metrics")
		defer wg.Done()

		if serveErr := s.Serve(); serveErr != nil {
			slog.
				With(slog.Any("err", serveErr)).
				Error("failed to serve metrics")
		}
	}(metricsServer)

	wg.Wait()

	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.Stop()
	cancel()
	slog.Info("[main] gracefully stopped")
}
