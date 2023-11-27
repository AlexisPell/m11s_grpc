package authHandlers

import (
	ssov1 "github.com/AlexisPell/m11s_grpc/protos/gen/go/sso"
	"github.com/AlexisPell/m11s_grpc/services/sso/internal/utils/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateLoginRequest(req *ssov1.LoginRequest) error {
	if errs := validator.Get().Var(req.GetEmail(), "email,required"); errs != nil {
		return status.Error(codes.InvalidArgument, "Email is not correct.")
	}
	if errs := validator.Get().Var(req.GetPassword(), "required,min=6"); errs != nil {
		return status.Error(codes.InvalidArgument, "Password must be at least 6 chars.")
	}
	if req.GetAppId() == 0 {
		return status.Error(codes.InvalidArgument, "AppId must be defined.")
	}
	return nil
}

func validateRegisterRequest(req *ssov1.RegisterRequest) error {
	if errs := validator.Get().Var(req.GetEmail(), "email,required"); errs != nil {
		return status.Error(codes.InvalidArgument, "Email is not correct.")
	}
	if errs := validator.Get().Var(req.GetPassword(), "required,min=6"); errs != nil {
		return status.Error(codes.InvalidArgument, "Password must be at least 6 chars.")
	}
	return nil
}

func validateIsAdminRequest(req *ssov1.IsAdminRequest) error {
	if errs := validator.Get().Var(req.GetUserId(), "required"); errs != nil {
		return status.Error(codes.InvalidArgument, "UserId is required")
	}
	return nil
}
