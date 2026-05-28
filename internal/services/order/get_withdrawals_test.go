package order

import (
	"context"
	"testing"
	"time"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_GetWithdrawals(t *testing.T) {
	var (
		mockRepository *MockRepository
		logger         zerolog.Logger
		userID         int
		repoResult     []models.DbWithdrawn
		processedTime  time.Time

		service *Service
	)

	setup := func(t *testing.T) {
		mockRepository = NewMockRepository(t)
		logger = zerolog.Nop()
		userID = 1
		processedTime = time.Now()

		repoResult = []models.DbWithdrawn{
			{
				ID:          1,
				UserID:      userID,
				ProcessedAt: processedTime,
				Sum:         100,
				Order:       "123",
			},
			{
				ID:          2,
				UserID:      userID,
				ProcessedAt: processedTime,
				Sum:         200,
				Order:       "321",
			},
		}

		service, _ = New(&logger, mockRepository, NewMockLoyaltyService(t))
	}

	t.Run(
		"Successfully", func(t *testing.T) {
			setup(t)

			mockRepository.EXPECT().GetWithdrawals(mock.Anything, userID).Return(repoResult, nil)

			res, err := service.GetWithdrawals(context.Background(), userID)

			require.NoError(t, err)
			assert.Equal(
				t, []models.ResponseWithdrawn{
					{
						Sum:         100,
						Order:       "123",
						ProcessedAt: processedTime,
					},
					{
						Sum:         200,
						Order:       "321",
						ProcessedAt: processedTime,
					},
				}, res,
			)
		},
	)

	t.Run(
		"Repository err", func(t *testing.T) {
			setup(t)

			testError := assert.AnError
			mockRepository.EXPECT().GetWithdrawals(mock.Anything, userID).Return(nil, testError)

			_, err := service.GetWithdrawals(context.Background(), userID)

			require.Error(t, err)
			assert.Equal(t, testError, err)
		},
	)

	t.Run(
		"Empty withdrawals list", func(t *testing.T) {
			setup(t)

			mockRepository.EXPECT().GetWithdrawals(mock.Anything, userID).Return([]models.DbWithdrawn{}, nil)

			result, err := service.GetWithdrawals(context.Background(), userID)

			require.NoError(t, err)
			assert.Empty(t, result)
		},
	)

	t.Run(
		"Multiple withdrawals successful retrieval", func(t *testing.T) {
			setup(t)

			withdrawals := []models.DbWithdrawn{
				{Sum: 10.0, Order: "12345678903", UserID: userID, ProcessedAt: processedTime},
				{Sum: 20.0, Order: "98765432107", UserID: userID, ProcessedAt: processedTime},
			}

			mockRepository.EXPECT().GetWithdrawals(mock.Anything, userID).Return(withdrawals, nil)

			result, err := service.GetWithdrawals(context.Background(), userID)

			require.NoError(t, err)
			assert.Len(t, result, 2)
			assert.Equal(t, 10.0, result[0].Sum)
			assert.Equal(t, "12345678903", result[0].Order)
		},
	)

	t.Run(
		"Database error on withdrawals list", func(t *testing.T) {
			setup(t)

			testError := assert.AnError
			mockRepository.EXPECT().GetWithdrawals(mock.Anything, userID).Return(nil, testError)

			result, err := service.GetWithdrawals(context.Background(), userID)

			require.Error(t, err)
			assert.Nil(t, result)
		},
	)
}
