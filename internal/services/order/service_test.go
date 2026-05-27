package order

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestNewUserService(t *testing.T) {
	t.Run(
		"Successfully", func(t *testing.T) {
			svc, err := New(&zerolog.Logger{}, NewMockRepository(t), NewMockLoyaltyService(t))

			require.NoError(t, err)
			require.NotNil(t, svc)
			require.NotNil(t, svc.logger)
			require.NotNil(t, svc.repository)
			require.NotNil(t, svc.loyaltyService)
		},
	)

	t.Run(
		"Failed, empty deps", func(t *testing.T) {
			logger := &zerolog.Logger{}
			repository := NewMockRepository(t)
			loyalty := NewMockLoyaltyService(t)

			tests := []struct {
				logger     *zerolog.Logger
				repository Repository
				loyalty    LoyaltyService
			}{
				{logger: nil, repository: repository, loyalty: loyalty},
				{logger: logger, repository: nil, loyalty: loyalty},
				{logger: logger, repository: repository, loyalty: nil},
			}

			for _, test := range tests {
				svc, err := New(test.logger, test.repository, test.loyalty)
				require.Error(t, err)
				require.Nil(t, svc)
			}
		},
	)
}
