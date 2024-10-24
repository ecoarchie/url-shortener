package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ecoarchie/url-shortener/cmd/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

	//TODO init logger: log/slog
	logger := setupLogger(cfg.Env)
	logger.Info("starting the app", slog.String("env", cfg.Env))
	logger.Debug("debug messages enabled")

	//TODO init storage: sqlite

	//TODO init router: chi, chi render

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
