package pg

import (
	"fmt"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Repository struct {
	logger *zerolog.Logger
	pool   *pgxpool.Pool
}

func New(logger *zerolog.Logger, pool *pgxpool.Pool) (*Repository, error) {
	if logger == nil {
		return nil, fmt.Errorf("%w: logger", errors.ErrEmptyDep)
	}
	if pool == nil {
		return nil, fmt.Errorf("%w: pgxpool", errors.ErrEmptyDep)
	}
	return &Repository{logger: logger, pool: pool}, nil
}
