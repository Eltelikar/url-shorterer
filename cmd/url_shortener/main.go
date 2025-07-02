package main

import (
	"log/slog"
	"os"
	"project_1/internal/config"
	"project_1/internal/storage/postgre"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	slog.SetDefault(log)

	slog.Info("Starting URL Shortener Service",
		slog.String("env", cfg.Env))

	slog.Debug("debug messages are enabled")

	// Initialize storage
	slog.Debug("Initializing storage")
	storage, err := postgre.New(config.GetStorageLink(cfg))
	if err != nil {
		log.Error("failed to init storage", slog.String("error", err.Error()))
		os.Exit(1)
	} else {
		slog.Debug("Storage initialized successfully")
	}

	err = storage.DeleteURL("yandextest")

	if err != nil {
		log.Error("failed to delete URL", slog.String("error", err.Error()))
		os.Exit(1)
	}
	_ = storage

	// TODO: init router: chi, "chi render"

	// TODO: run server:
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
