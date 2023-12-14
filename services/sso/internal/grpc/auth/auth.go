package auth

import (
	"context"
	"errors"
	"fmt"
	ssov1 "github.com/AlexisPell/m11s_grpc/protos/gen/go/sso"
	"github.com/AlexisPell/m11s_grpc/services/sso/internal/services/auth"
	"github.com/AlexisPell/m11s_grpc/services/sso/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type authHandlers struct {
	log                           *slog.Logger
	authService                   auth.AuthService
	ssov1.UnimplementedAuthServer // temporarily plug for unimplemented handlers
}

// Register To register our auth service on grpc server
func Register(gRPC *grpc.Server, log *slog.Logger, authService auth.AuthService) {
	ssov1.RegisterAuthServer(gRPC, &authHandlers{log: log, authService: authService})
}

func (h *authHandlers) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	// validate request
	if err := validateLoginRequest(req); err != nil {
		return nil, err
	}

	// pass to service
	token, err := h.authService.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "Invalid credentials")
		}
		return nil, status.Error(codes.Internal, "Internal error occurred")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (h *authHandlers) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	// validate request
	if err := validateRegisterRequest(req); err != nil {
		return nil, err
	}

	// pass to service
	userId, err := h.authService.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		fmt.Println("WE ARE HERE", err)
		if errors.Is(err, storage.ErrUserExists) {
			fmt.Println("WE ARE HERE 2", err)
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "Internal error occurred")
	}

	return &ssov1.RegisterResponse{
		UserId: int64(userId),
	}, nil
}

func (h *authHandlers) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	// validate request
	if err := validateIsAdminRequest(req); err != nil {
		return nil, err
	}

	// pass to service
	isAdmin, err := h.authService.IsAdmin(ctx, int(req.GetUserId()))
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "Internal error occurred")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
