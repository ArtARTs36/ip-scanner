package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"runtime/debug"
	"strings"

	goMetrics "github.com/artarts36/go-metrics"
	"github.com/artarts36/ip-scanner/internal/domain"
	"github.com/artarts36/ip-scanner/internal/infrastructure/repository"
	"github.com/artarts36/ip-scanner/internal/infrastructure/storage"
	"github.com/artarts36/ip-scanner/internal/metrics"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type App struct {
	gRPCServer *grpc.Server
	port       int
	cfg        *Config

	infrastructure struct {
		repositories struct {
			ipRepository domain.IPRepository
		}

		storage *storage.AutoUpdater
	}

	metrics *metrics.Metrics
}

// NewApp creates new gRPC server app.
func NewApp(
	ctx context.Context,
	cfg *Config,
) (*App, error) {
	app := &App{
		port: cfg.GRPC.Port,
		cfg:  cfg,
	}

	app.metrics = metrics.NewMetrics(goMetrics.NewDefaultRegistry(cfg.Metrics.ToGoMetrics()))

	if err := app.initRepositories(ctx); err != nil {
		return nil, fmt.Errorf("init repositories: %w", err)
	}

	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			// logging.StartCall, logging.FinishCall,
			logging.PayloadReceived, logging.PayloadSent,
		),
		// Add any other option (check functions starting with logging.With).
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) error {
			slog.
				With(slog.Any("stack", getStackTrace())).
				Error("[app] recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	app.gRPCServer = grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(interceptorLogger(), loggingOpts...),
	))

	app.registerServices()

	if cfg.GRPC.UseReflection {
		reflection.Register(app.gRPCServer)
	}

	return app, nil
}

// interceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func interceptorLogger() logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, _ ...any) {
		switch lvl {
		case logging.LevelDebug:
			slog.DebugContext(ctx, msg)
		case logging.LevelInfo:
			slog.InfoContext(ctx, msg)
		case logging.LevelWarn:
			slog.WarnContext(ctx, msg)
		case logging.LevelError:
			slog.ErrorContext(ctx, msg)
		}
	})
}

// Run runs gRPC server.
func (app *App) Run(ctx context.Context) error {
	listenConfig := net.ListenConfig{}

	l, err := listenConfig.Listen(ctx, "tcp", fmt.Sprintf(":%d", app.port))
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	slog.Info("[grpc] grpc server started")

	if serveErr := app.gRPCServer.Serve(l); serveErr != nil {
		return fmt.Errorf("serve: %w", serveErr)
	}

	return nil
}

// Stop stops gRPC server.
func (app *App) Stop() {
	slog.Info("[grpc] stopping gRPC server")

	app.gRPCServer.GracefulStop()
}

func (app *App) initRepositories(ctx context.Context) error {
	slog.Debug("[container] loading maxminddb")

	if err := app.initStorage(); err != nil {
		return fmt.Errorf("init storage: %w", err)
	}

	mmdb, err := app.infrastructure.storage.Load(ctx)
	if err != nil {
		return fmt.Errorf("load db: %w", err)
	}

	app.metrics.MMDB.Scrap(mmdb.ActualTime, mmdb.DB.Metadata.DatabaseType, mmdb.DB.Metadata.Languages, mmdb.Size)

	if app.cfg.Storage.AutoUpdate.Enabled() {
		go func() {
			app.infrastructure.storage.Run()
		}()
	}

	slog.Debug("[container] maxminddb loaded")

	app.infrastructure.repositories.ipRepository = repository.NewIPRepository(mmdb)

	return nil
}

func (app *App) initStorage() error {
	storages, err := storage.Create(app.cfg.Storage)
	if err != nil {
		return err
	}

	app.infrastructure.storage = storage.NewAutoUpdater(storages, app.cfg.Storage, app.metrics)

	return nil
}

func getStackTrace() []string {
	stack := strings.ReplaceAll(string(debug.Stack()), "\t", "")
	stackRows := strings.Split(stack, "\n")
	return stackRows
}
