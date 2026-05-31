// Package order содержит обработчики http запросов связанные с обработкой заказов
package order

import (
	"context"
	"fmt"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

// Service интерфейс сервиса с доменной логикой
type Service interface {
	Upload(ctx context.Context, userID int, orderNum string) error
	GetList(ctx context.Context, userID int) ([]models.RequestOrder, error)
	Withdraw(ctx context.Context, userID int, sum float64, orderNum string) error
	GetBalance(ctx context.Context, userID int) (models.ResponseBalance, error)
	GetWithdrawals(ctx context.Context, userID int) ([]models.ResponseWithdrawn, error)
}

// AuthService интерфейс сервиса получения данных от системы расчета баллов
type AuthService interface {
	GetUserID(ctx context.Context) (int, error)
}

// Handler структура обработчика запросов
type Handler struct {
	logger      *zerolog.Logger
	service     Service
	validate    *validator.Validate
	authService AuthService
}

// New конструктор обработчика запросов
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
