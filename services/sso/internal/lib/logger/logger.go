package logger

import (
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
	"sync"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

var log *slog.Logger
var once sync.Once

func SetupLogger(env string) *slog.Logger {
	once.Do(func() {
		switch env {
		case envLocal:
			log = slog.New(
				tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}),
			)
		case envDev:
			log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		case envProd:
			log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		}
	})

	return log
}

func Get() *slog.Logger {
	return log
}
