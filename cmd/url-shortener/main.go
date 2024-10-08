package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/IceTweak/url-shortener/internal/config"
	"github.com/IceTweak/url-shortener/internal/http-server/handlers/url/save"
	mwLogger "github.com/IceTweak/url-shortener/internal/http-server/middleware/logger"
	"github.com/IceTweak/url-shortener/internal/lib/logger/sl"
	"github.com/IceTweak/url-shortener/internal/lib/logger/sl/handlers/slogpretty"
	"github.com/IceTweak/url-shortener/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	router := chi.NewRouter()

	// middlewares
	router.Use(middleware.RequestID)
	// chi's logger middleware
	router.Use(middleware.Logger)
	// own implemented logger middleware
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// handlers
	router.Post("/url", save.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	// init server
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// start server
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
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
	// If env config is invalid, set prod settings by default due to security
	default:
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
