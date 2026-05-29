// Package middlewares кастомные middleware приложения + служба для получения данных пользователя
package middlewares

import (
	"fmt"

	domainErrors "github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/rs/zerolog"
)

type Config interface {
	GetSalt() []byte
}

// Middleware структура для middlewares
type Middleware struct {
	logger *zerolog.Logger
	cfg    Config
}

// New конструктор
func New(logger *zerolog.Logger, cfg Config) (*Middleware, error) {
	if logger == nil {
		return nil, fmt.Errorf("%w: logger", domainErrors.ErrEmptyDep)
	}
	if cfg == nil {
		return nil, fmt.Errorf("%w: config", domainErrors.ErrEmptyDep)
	}
	return &Middleware{logger: logger, cfg: cfg}, nil
}
