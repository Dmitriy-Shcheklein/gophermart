package user

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestNewUserService(t *testing.T) {
	t.Run(
		"Successfully", func(t *testing.T) {
			svc, err := New(&zerolog.Logger{}, NewMockRepository(t), NewMockConfig(t))

			require.NoError(t, err)
			require.NotNil(t, svc)
			require.NotNil(t, svc.logger)
			require.NotNil(t, svc.repository)
			require.NotNil(t, svc.cfg)
		},
	)

	t.Run(
		"Failed, empty deps", func(t *testing.T) {
			logger := &zerolog.Logger{}
			repository := NewMockRepository(t)
			cfg := NewMockConfig(t)

			tests := []struct {
				logger     *zerolog.Logger
				repository Repository
				cfg        Config
			}{
				{logger: nil, repository: repository, cfg: cfg},
				{logger: logger, repository: nil, cfg: cfg},
				{logger: logger, repository: repository, cfg: nil},
			}

			for _, test := range tests {
				svc, err := New(test.logger, test.repository, test.cfg)
				require.Error(t, err)
				require.Nil(t, svc)
			}
		},
	)
}
