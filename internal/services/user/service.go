// Package user пакет по работе с пользователем
package user

import (
	"context"
	"fmt"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

// Repository интерфейс описывающий методы репозитория
type Repository interface {
	CreateUser(ctx context.Context, login string, password string, balance models.DbBalance) (models.DbUser, error)
	GetUserByLogin(ctx context.Context, login string) (*models.DbUser, error)
}

// Config интерфейс описывающий конфигурацию приложения
type Config interface {
	GetSalt() []byte
}

// Service структура сервиса
type Service struct {
	logger     *zerolog.Logger
	repository Repository
	cfg        Config
}

// New конструктор сервиса
func New(logger *zerolog.Logger, repository Repository, cfg Config) (*Service, error) {
	if logger == nil {
		return nil, fmt.Errorf("%w: logger", errors.ErrEmptyDep)
	}
	if repository == nil {
		return nil, fmt.Errorf("%w: repository", errors.ErrEmptyDep)
	}
	if cfg == nil {
		return nil, fmt.Errorf("%w: config", errors.ErrEmptyDep)
	}
	return &Service{logger: logger, repository: repository, cfg: cfg}, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}
