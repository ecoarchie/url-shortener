package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/ecoarchie/url-shortener/internal/config"
	deleter "github.com/ecoarchie/url-shortener/internal/http-server/handlers/url/deleter"
	"github.com/ecoarchie/url-shortener/internal/http-server/handlers/url/redirect"
	"github.com/ecoarchie/url-shortener/internal/http-server/handlers/url/save"
	mvlogger "github.com/ecoarchie/url-shortener/internal/http-server/middleware/logger"
	"github.com/ecoarchie/url-shortener/internal/lib/logger/slg"
	"github.com/ecoarchie/url-shortener/internal/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// init config
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := config.MustLoad()

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

	// init router: chi, chi render
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger) // chi logger
	router.Use(mvlogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat) // allows working with url vars

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password, // pair for 1 user-password
		}))

		r.Post("/", save.New(logger, storage))
		r.Delete("/{alias}", deleter.New(logger, storage))
	})

	router.Get("/{alias}", redirect.New(logger, storage))

	logger.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// run server
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("failed to start server")
	}
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
