package app

import (
	grpcapp "github.com/AlexisPell/m11s_grpc/services/sso/internal/app/grpc"
	"github.com/AlexisPell/m11s_grpc/services/sso/internal/services/auth"
	"github.com/AlexisPell/m11s_grpc/services/sso/internal/storage/sqlite"
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
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	// init auth service
	authServ := auth.New(log, storage, storage, tokenTTL)

	// init grpc app
	grpcApp := grpcapp.New(log, authServ, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
