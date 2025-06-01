// @title API системы электронного дневника
// @version 1.0
// @description Полная документация API для электронного дневника
// @contact.name API Support
// @license.name MIT
// @host localhost:8080

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log/slog"
	"os"
	"os/signal"
	"service/internal/config"
	"service/internal/http-server/handler"
	"service/internal/lib/logger/handlers/slogpretty"
	"service/internal/lib/logger/sl"
	"service/internal/storage/mysql"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting edu-helper", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := mysql.New(cfg.SQLPath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	srv, err := handler.NewServer(log, cfg, storage)
	if err != nil {
		log.Error("failed to init http server", sl.Err(err))
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
