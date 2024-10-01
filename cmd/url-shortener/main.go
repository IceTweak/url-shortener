package main

import (
	"log/slog"
	"os"

	"github.com/IceTweak/url-shortener/internal/config"
	"github.com/IceTweak/url-shortener/internal/lib/logger/sl"
	"github.com/IceTweak/url-shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("🚀 Starting url-shortener...", slog.String("env", cfg.Env))
	log.Debug("🛠️ Debug level enabled!")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("🆘 Failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	// TODO: init router - chi, chi render

	// TODO: run server
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
