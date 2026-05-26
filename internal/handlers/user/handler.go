package user

import (
	"context"
	"fmt"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type Service interface {
	Register(ctx context.Context, login string, password string) error
	Auth(ctx context.Context, login string, password string) (string, error)
}

type Handler struct {
	logger   *zerolog.Logger
	service  Service
	validate *validator.Validate
}

func New(logger *zerolog.Logger, service Service) (*Handler, error) {
	if logger == nil {
		return nil, fmt.Errorf("%w: logger", errors.ErrEmptyDep)
	}
	if service == nil {
		return nil, fmt.Errorf("%w: service", errors.ErrEmptyDep)
	}
	return &Handler{
		logger: logger, service: service, validate: validator.New(validator.WithRequiredStructEnabled()),
	}, nil
}
