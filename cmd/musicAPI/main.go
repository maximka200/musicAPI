package main

import (
	"context"
	"fmt"
	"log/slog"
	"musicAPI/internal/config"
	"musicAPI/internal/repository/psql"
	"musicAPI/internal/services"
	"musicAPI/internal/transport/client/musicInfo"
	"musicAPI/internal/transport/handlers"
	"musicAPI/internal/transport/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustReadConfig()
	log := initLogger(cfg.Env)
	log.Info("starting application", slog.Any("config", cfg))

	serv := server.Server{}
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	client := musicInfo.NewMusicInfo(cfg.ApiAddress, cfg.Timeout)
	repos := psql.MustNewDB(&cfg)
	service := services.NewService(client, repos)
	handler := handlers.NewHandler(log, service, ctx)

	go func() {
		if err := serv.Run(cfg, handler.InitRouter()); err != nil {
			log.Error(fmt.Sprintf("cannot run server: %s", err))
			panic("cannot run server")
		}
	}()

	// shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	if err := serv.Shutdown(ctx); err != nil {
		log.Error(fmt.Sprintf("an error occurred while executing graceful shutdown: %s", err))
	}
}

func initLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "dev":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
