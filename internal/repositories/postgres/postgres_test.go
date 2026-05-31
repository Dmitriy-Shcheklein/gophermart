package postgres

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestNewUserRepository(t *testing.T) {
	t.Run(
		"Successfully", func(t *testing.T) {
			repo, err := New(&zerolog.Logger{}, &pgxpool.Pool{})

			require.NoError(t, err)
			require.NotNil(t, repo)
		},
	)

	t.Run(
		"Failed, empty deps", func(t *testing.T) {
			logger := &zerolog.Logger{}
			pool := &pgxpool.Pool{}

			tests := []struct {
				logger *zerolog.Logger
				pool   *pgxpool.Pool
			}{
				{logger: nil, pool: pool},
				{logger: logger, pool: nil},
			}

			for _, test := range tests {
				repo, err := New(test.logger, test.pool)
				require.Error(t, err)
				require.Nil(t, repo)
			}
		},
	)
}
