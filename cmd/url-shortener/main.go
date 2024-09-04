package main

import (
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"

	mwLogger "url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/lib/logger/handlers/slogpretty"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "development"
	envProd  = "production"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting url-shortener")
	log.Debug("debug msg enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	id, err := storage.SaveURL("https://google.com", "google")
	if err != nil {
		log.Error("Failed to save url", sl.Err(err))
		os.Exit(1)
	}

	log.Info("Saved URL", slog.Int64("id", id))

	err = storage.DeleteURL("google")
	if err != nil {
		log.Error("Failed to delete url", sl.Err(err))
		os.Exit(1)
	}

	log.Info("Deleted URL with alias 'google'")

	_ = storage

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
}

func setupLogger(env string) *slog.Logger {
	var opts slogpretty.PrettyHandlerOptions

	switch env {
	case envLocal:
		opts = slogpretty.PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}

	case envDev:
		opts = slogpretty.PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}

	case envProd:
		opts = slogpretty.PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		}

	default: // Если окружение не распознано, используем настройки для продакшена по умолчанию
		opts = slogpretty.PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		}
	}

	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
