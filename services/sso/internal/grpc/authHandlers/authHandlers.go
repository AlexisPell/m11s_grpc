package authHandlers

import (
	"context"
	ssov1 "github.com/AlexisPell/m11s_grpc/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type AuthService interface {
	Login(ctx context.Context, email, password string, appId int) (token string, err error)
	Register(ctx context.Context, email, password string) (userId int, err error)
	IsAdmin(ctx context.Context, userId int) (isAdmin bool, err error)
}

type authHandlers struct {
	log                           *slog.Logger
	authService                   AuthService
	ssov1.UnimplementedAuthServer // temporarily plug for unimplemented handlers
}

// Register To register our authService service on grpc server
func Register(gRPC *grpc.Server, log *slog.Logger, authService AuthService) { // authService AuthService
	ssov1.RegisterAuthServer(gRPC, &authHandlers{log: log, authService: authService}) // authService: authService,
}

func (h *authHandlers) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	// validate request
	if err := validateLoginRequest(req); err != nil {
		return nil, err
	}

	// pass to service
	//token, err := h.authService.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	//if err != nil {
	//	// TODO: handle error
	//	return nil, status.Error(codes.Internal, "Internal error occurred")
	//}

	return &ssov1.LoginResponse{
		Token: "token",
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
		// TODO: handle error
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
		// TODO: handle error
		return nil, status.Error(codes.Internal, "Internal error occurred")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
