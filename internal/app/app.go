package app

import (
	"context"
	"github.com/ensiouel/static/internal/config"
	"github.com/ensiouel/static/internal/domain/static/service"
	"github.com/ensiouel/static/internal/domain/static/storage"
	static "github.com/ensiouel/static/internal/domain/static/transport/grpc/v1"
	"github.com/ensiouel/static/internal/transport/grpc"
	"golang.org/x/exp/slog"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
}

func New() *App {
	return &App{}
}

func (app *App) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	conf := config.New()

	var level slog.Level
	level.UnmarshalText([]byte(conf.LogLevel))

	logger := slog.New(slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: level,
		},
	))

	staticStorage, err := storage.NewStaticStorage(conf.Static.StorageRoot)
	if err != nil {
		logger.Error("failed to initialize static storage",
			slog.Any("error", err),
			slog.String("root", conf.Static.StorageRoot),
		)
		return
	}

	staticService := service.NewStaticService(logger, staticStorage)
	staticServer := static.NewStaticServer(staticService, conf.Static.MaxFileSize)

	server := grpc.New(logger).
		Register(staticServer)

	go func() {
		logger.Info("starting server", slog.Any("addr", conf.GRPC.Addr))
		err := server.Run(conf.GRPC.Addr)
		if err != nil {
			logger.Error("failed to listen", slog.Any("error", err))
			return
		}
	}()

	<-ctx.Done()

	logger.Info("shutting down server")

	server.Stop()
}
