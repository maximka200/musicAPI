package main

import (
	"log/slog"
	"musicAPI/internal/config"
	"os"
)

func main() {
	cfg := config.MustReadConfig()
	log := initLogger(cfg.Env)

	/* go func() {
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
	} */
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
