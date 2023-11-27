package main

import (
	"github.com/AlexisPell/m11s_grpc/services/sso/internal/app"
	"github.com/AlexisPell/m11s_grpc/services/sso/internal/config"
	"github.com/AlexisPell/m11s_grpc/services/sso/internal/lib/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// Short description
// Main - entrypoint
// App - application wrapper
// App.GRPCServer - application for grpc handling
// TODO: App.Database

func main() {
	// get config
	cfg := config.MustLoad()

	// logger
	log := logger.SetupLogger(cfg.Env)

	log.Info("Starting application. config: ", slog.Any("cfg", cfg))

	// app initialization
	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	// launch grpc-server
	go application.GRPCServer.MustRun()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	// listen for 2 system calls and if they occur - place this event into the stop channel
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// Main goroutine is blocked until smth is written to stop channel
	sig := <-stop
	log.Warn("Shutting down application.", slog.String("signal", sig.String()))
	application.GRPCServer.Stop()
	// TODO: graceful shutdown for database connection
}
