package order

import (
	"context"
	"fmt"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
	"github.com/rs/zerolog"
)

type Repository interface {
	CreateOrder(ctx context.Context, userID int, orderNum string) error
	GetOrderByNum(ctx context.Context, orderNum string) (models.DbOrder, error)
	GetByUserId(ctx context.Context, userID int) ([]models.DbOrder, error)
	GetUserBalance(ctx context.Context, useID int) (models.DbBalance, error)
	Withdraw(ctx context.Context, balance models.DbBalance, withdrawn models.DbWithdrawn) error
	GetProcessingOrders(ctx context.Context) ([]models.DbOrder, error)
	UpdateOrder(ctx context.Context, order models.DbOrder) error
	GetWithdrawals(ctx context.Context, userID int) ([]models.DbWithdrawn, error)
}

type LoyaltyService interface {
	GetOrder(ctx context.Context, orderNum string) (models.LoyaltyOrderData, error)
}

type Service struct {
	logger         *zerolog.Logger
	repository     Repository
	loyaltyService LoyaltyService
}

func New(logger *zerolog.Logger, repository Repository, loyalty LoyaltyService) (*Service, error) {
	if logger == nil {
		return nil, fmt.Errorf("%w: logger", errors.ErrEmptyDep)
	}
	if repository == nil {
		return nil, fmt.Errorf("%w: repository", errors.ErrEmptyDep)
	}
	if loyalty == nil {
		return nil, fmt.Errorf("%w: loyaltyService", errors.ErrEmptyDep)
	}
	return &Service{logger: logger, repository: repository, loyaltyService: loyalty}, nil
}
