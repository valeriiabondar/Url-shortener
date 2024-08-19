package main

import (
	"log/slog"
	"os"
	"urlShortener/internal/config"
	"urlShortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setUpLogger(cfg.Env)
	log.Info("starting url shortener", slog.String("env", cfg.Env))

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("could not init storage", err)
		os.Exit(1)
	}

	id, err := storage.SaveUrl("https://www.google.com", "google")
	if err != nil {
		log.Error("could not save url", err)
	} else {
		log.Info("url saved", slog.Int64("id", id))
	}

	urlToGet, err := storage.GetUrl("fkf")
	if err != nil {
		log.Error("could not get url", err)
	} else {
		log.Info("url retrieved", slog.String("url", urlToGet))
	}

	err = storage.DeleteUrl("google")
	if err != nil {
		log.Error("could not delete url", err)
	} else {
		log.Info("url deleted")
	}
}

func setUpLogger(env string) *slog.Logger {
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
