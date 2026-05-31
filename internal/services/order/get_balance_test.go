package order

import (
	"context"
	"testing"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_GetBalance(t *testing.T) {
	var (
		mockRepository *MockRepository
		logger         zerolog.Logger
		userID         int
		repoResult     models.DbBalance

		service *Service
	)

	setup := func(t *testing.T) {
		mockRepository = NewMockRepository(t)
		logger = zerolog.Nop()
		userID = 1
		repoResult = models.DbBalance{
			ID:        1,
			Current:   10,
			Withdrawn: 100,
			UserID:    1,
		}

		service, _ = New(&logger, mockRepository, NewMockLoyaltyService(t))
	}

	t.Run(
		"Successfully", func(t *testing.T) {
			setup(t)

			mockRepository.EXPECT().GetUserBalance(mock.Anything, userID).Return(repoResult, nil)

			res, err := service.GetBalance(context.Background(), userID)

			require.NoError(t, err)
			assert.Equal(
				t, models.ResponseBalance{
					Current:   10,
					Withdrawn: 100,
				}, res,
			)
		},
	)

	t.Run(
		"Repository error", func(t *testing.T) {
			setup(t)

			testError := assert.AnError
			mockRepository.EXPECT().GetUserBalance(mock.Anything, userID).Return(repoResult, testError)

			res, err := service.GetBalance(context.Background(), userID)

			require.ErrorIs(t, err, testError)
			assert.Empty(t, res)
		},
	)

	t.Run(
		"Empty response from database", func(t *testing.T) {
			setup(t)

			emptyBalance := models.DbBalance{Current: 0.0, Withdrawn: 0.0, UserID: userID, ID: 1}
			mockRepository.EXPECT().GetUserBalance(mock.Anything, userID).Return(emptyBalance, nil)

			result, err := service.GetBalance(context.Background(), userID)

			require.NoError(t, err)
			assert.Equal(t, 0.0, result.Current)
		},
	)

	t.Run(
		"Large balance values", func(t *testing.T) {
			setup(t)

			largeBalance := models.DbBalance{Current: 999999.99, Withdrawn: 888888.88, UserID: userID, ID: 1}
			mockRepository.EXPECT().GetUserBalance(mock.Anything, userID).Return(largeBalance, nil)

			result, err := service.GetBalance(context.Background(), userID)

			require.NoError(t, err)
			assert.Equal(t, 999999.99, result.Current)
			assert.Equal(t, 888888.88, result.Withdrawn)
		},
	)
}
