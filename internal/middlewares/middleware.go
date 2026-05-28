// Package middlewares кастомные middleware приложения + служба для получения данных пользователя
package middlewares

import (
	"fmt"

	domainErrors "github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/rs/zerolog"
)

// Middleware структура для middlewares
type Middleware struct {
	logger *zerolog.Logger
}

// New конструктор
func New(logger *zerolog.Logger) (*Middleware, error) {
	if logger == nil {
		return nil, fmt.Errorf("%w: logger", domainErrors.ErrEmptyDep)
	}
	return &Middleware{logger: logger}, nil
}
