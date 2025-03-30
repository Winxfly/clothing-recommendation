package main

import (
	"clothing-recommendation/internal/config"
	"clothing-recommendation/internal/http-server/handlers/geocode"
	"clothing-recommendation/internal/http-server/handlers/recommendation"
	"clothing-recommendation/internal/http-server/middleware/cors"
	"clothing-recommendation/internal/http-server/middleware/logger"
	loggerSlog "clothing-recommendation/internal/lib/logger/slog"
	"clothing-recommendation/internal/storage/postgresql"
	"clothing-recommendation/internal/weather"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	os.Setenv("CONFIG_PATH", "./config/local.yaml")

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("starting clothing-recommendation", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := postgresql.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", loggerSlog.Err(err))
		os.Exit(1)
	}

	geoClient := geocode.New(cfg.Geocoding)
	weatherClient := weather.New(cfg.Weather)

	mux := http.NewServeMux()
	mux.Handle("GET /geocode", geocode.Handler(log, geoClient))
	mux.Handle("POST /recommend", recommendation.New(log, storage, weatherClient))

	handler := logger.New(log)(
		cors.Middleware(
			mux,
		),
	)

	log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))
	if err := http.ListenAndServe(cfg.HTTPServer.Address, handler); err != nil {
		log.Error("failed to start server", loggerSlog.Err(err))
		os.Exit(1)
	}

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
