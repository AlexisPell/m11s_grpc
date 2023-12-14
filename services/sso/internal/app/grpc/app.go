package grpc_app

import (
	"fmt"
	authHandlers "github.com/AlexisPell/m11s_grpc/services/sso/internal/grpc/auth"
	authService "github.com/AlexisPell/m11s_grpc/services/sso/internal/services/auth"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, authService authService.AuthService, port int) *App {
	gRPCServer := grpc.NewServer()

	// define services for handlers here
	//auth

	authHandlers.Register(gRPCServer, log, authService)

	return &App{
		gRPCServer: gRPCServer,
		log:        log,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpc_app.Run"
	log := a.log.With(slog.String("op", op))

	// Listen to TCP
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("GRPC Server is running ", slog.String("addr", l.Addr().String()))

	// Start GRPC on TCP
	err = a.gRPCServer.Serve(l)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpc_app.Stop"

	a.log.With(slog.String("op", op)).Info("Stopping GRPC Server. Port: ", a.port)

	a.gRPCServer.GracefulStop()
}
