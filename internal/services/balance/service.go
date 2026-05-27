package order

import (
	"context"
	"fmt"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/rs/zerolog"
)

type Repository interface {
	SetBalance(ctx context.Context, userID int, balance float64) error
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
