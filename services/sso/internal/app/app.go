package app

import (
	grpcapp "github.com/AlexisPell/m11s_grpc/services/sso/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	// initialize storage

	// init authService service

	// init grpc app
	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
