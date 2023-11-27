package authService

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexisPell/m11s_grpc/services/sso/internal/domain/models"
	"github.com/AlexisPell/m11s_grpc/services/sso/internal/lib/jwt"
	"github.com/AlexisPell/m11s_grpc/services/sso/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("Invalid credentials")
	ErrUserExists         = errors.New("User exists")
)

type authService struct {
	log          *slog.Logger
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type AuthService interface {
	Login(ctx context.Context, email, password string, appId int) (token string, err error)
	Register(ctx context.Context, email, password string) (userId int, err error)
	IsAdmin(ctx context.Context, userId int) (isAdmin bool, err error)
}

type UserProvider interface {
	SaveUser(ctx context.Context, email string, passHash string) (uid int64, err error)
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, id int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appId int) (models.App, error)
}

func New(log *slog.Logger, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *authService {
	return &authService{
		log:          log,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

// Login checks if the user exists and credentials are valid
// if not, returns error
func (s *authService) Login(ctx context.Context, email string, password string, appId int) (token string, err error) {
	const op = "authService.Login"

	log := s.log.With(slog.String("op", op))
	log.Info("Login attempt", slog.String("email", email))

	user, err := s.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("User not found" + err.Error())

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		log.Warn("Failed to get user: ", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Error("Invalid credentials: ", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := s.appProvider.App(ctx, appId)

	token, err = jwt.NewToken(user, app, s.tokenTTL)
	if err != nil {
		log.Error("Failed to create token", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}

// Login checks if the user exists and credentials are valid
// if not, returns error
func (s *authService) Register(ctx context.Context, email string, password string) (int64, error) {
	const op = "authService.Register"

	log := s.log.With(slog.String("op", op))

	log.Info("User register attempt", slog.String("email", email))

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to hash password")
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := s.userProvider.SaveUser(ctx, email, string(passHash))
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("User already exists", slog.String("error", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		log.Error("Failed to save user")
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("User registered", slog.Int64("userId", id))

	return id, nil
}

// Login checks if the user exists and credentials are valid
// if not, returns error
func (s *authService) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const op = "authService.IsAdmin"

	log := s.log.With(slog.String("op", op))
	log.Info("Check if user is admin")

	isAdmin, err := s.userProvider.IsAdmin(ctx, userId)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("User not found", slog.String("error", err.Error()))
		}
		log.Warn("Failed to check if user is admin", slog.String("error", err.Error()))
		return false, err
	}
	return isAdmin, nil
}
