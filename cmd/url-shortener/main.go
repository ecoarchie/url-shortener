package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ecoarchie/url-shortener/cmd/internal/config"
	mvlogger "github.com/ecoarchie/url-shortener/cmd/internal/http-server/middleware/logger"
	"github.com/ecoarchie/url-shortener/cmd/internal/lib/logger/slg"
	"github.com/ecoarchie/url-shortener/cmd/internal/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// init config
	cfg := config.MustLoad()

	fmt.Println(cfg)

	// init logger: log/slog
	logger := setupLogger(cfg.Env)
	logger.Info("starting the app", slog.String("env", cfg.Env))
	logger.Debug("debug messages enabled")

	// init storage: sqlite
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Error("failed to initialize storage", slg.Err(err))
		os.Exit(1)
	}

	//TODO init router: chi, chi render
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger) // chi logger
	router.Use(mvlogger.New(logger))

	//TODO run server
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}
