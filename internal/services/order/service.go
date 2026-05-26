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
}

type Service struct {
	logger     *zerolog.Logger
	repository Repository
}

func New(logger *zerolog.Logger, repository Repository) (*Service, error) {
	if logger == nil {
		return nil, fmt.Errorf("%w: logger", errors.ErrEmptyDep)
	}
	if repository == nil {
		return nil, fmt.Errorf("%w: repository", errors.ErrEmptyDep)
	}
	return &Service{logger: logger, repository: repository}, nil
}
