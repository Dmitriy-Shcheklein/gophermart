package order

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestNewUserService(t *testing.T) {
	t.Run(
		"Successfully", func(t *testing.T) {
			svc, err := New(&zerolog.Logger{}, NewMockRepository(t))

			require.NoError(t, err)
			require.NotNil(t, svc)
		},
	)

	t.Run(
		"Failed, empty deps", func(t *testing.T) {
			logger := &zerolog.Logger{}
			repository := NewMockRepository(t)

			tests := []struct {
				logger     *zerolog.Logger
				repository Repository
			}{
				{logger: nil, repository: repository},
				{logger: logger, repository: nil},
			}

			for _, test := range tests {
				svc, err := New(test.logger, nil)
				require.Error(t, err)
				require.Nil(t, svc)
			}
		},
	)
}
