package user

import (
	"context"
	"fmt"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	CreateUser(ctx context.Context, login string, password string) (models.DbUser, error)
	GetUserByLogin(ctx context.Context, login string) (*models.DbUser, error)
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

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}
