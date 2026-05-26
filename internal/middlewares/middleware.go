package middlewares

import (
	"fmt"

	domainErrors "github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/rs/zerolog"
)

type Middleware struct {
	logger *zerolog.Logger
}

func New(logger *zerolog.Logger) (*Middleware, error) {
	if logger == nil {
		return nil, fmt.Errorf("%w: logger", domainErrors.ErrEmptyDep)
	}
	return &Middleware{logger: logger}, nil
}
