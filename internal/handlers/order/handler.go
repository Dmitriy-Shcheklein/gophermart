package user

import (
	"context"
	"fmt"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type Service interface {
	Upload(ctx context.Context, userID int, orderNum string) error
	GetList(ctx context.Context, userID int) ([]models.RequestOrder, error)
	Withdraw(ctx context.Context, userID int, sum float64, orderNum string) error
	GetBalance(ctx context.Context, userID int) (models.ResponseBalance, error)
	GetWithdrawals(ctx context.Context, userID int) ([]models.ResponseWithdrawn, error)
}

type AuthService interface {
	GetUserID(ctx context.Context) (int, error)
}

type Handler struct {
	logger      *zerolog.Logger
	service     Service
	validate    *validator.Validate
	authService AuthService
}

func New(logger *zerolog.Logger, service Service, authService AuthService) (*Handler, error) {
	if logger == nil {
		return nil, fmt.Errorf("%w: logger", errors.ErrEmptyDep)
	}
	if service == nil {
		return nil, fmt.Errorf("%w: service", errors.ErrEmptyDep)
	}
	if authService == nil {
		return nil, fmt.Errorf("%w: authService", errors.ErrEmptyDep)
	}
	return &Handler{
		logger: logger, service: service, validate: validator.New(validator.WithRequiredStructEnabled()),
		authService: authService,
	}, nil
}
