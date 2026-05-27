package order

import (
	"context"
	"testing"

	domainErrors "github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_Withdraw(t *testing.T) {
	var (
		mockRepository *MockRepository
		logger         zerolog.Logger
		userID         int
		orderNum       string
		sum            float64
		balance        models.DbBalance

		service *Service
	)

	setup := func(t *testing.T) {
		mockRepository = NewMockRepository(t)
		logger = zerolog.Nop()
		userID = 1
		orderNum = "0"
		sum = 100.22
		balance = models.DbBalance{
			ID:        1,
			Current:   sum + 0.01,
			Withdrawn: 1000,
			UserID:    userID,
		}

		service, _ = New(&logger, mockRepository, NewMockLoyaltyService(t))
	}

	t.Run(
		"Successfully", func(t *testing.T) {
			setup(t)

			mockRepository.EXPECT().GetUserBalance(mock.Anything, userID).Return(balance, nil)
			mockRepository.EXPECT().Withdraw(
				mock.Anything,
				models.DbBalance{UserID: balance.UserID, ID: balance.ID, Current: 0.01, Withdrawn: 1100.22},
				models.DbWithdrawn{Order: orderNum, Sum: sum, UserID: userID},
			).Return(nil)

			err := service.Withdraw(context.Background(), userID, sum, orderNum)

			require.NoError(t, err)
		},
	)

	t.Run(
		"Error while getting balance", func(t *testing.T) {
			setup(t)

			testError := assert.AnError
			mockRepository.EXPECT().GetUserBalance(mock.Anything, mock.Anything).Return(balance, testError)

			err := service.Withdraw(context.Background(), userID, sum, orderNum)

			require.Error(t, err)
			assert.Equal(t, testError, err)
		},
	)

	t.Run(
		"Balance less 0", func(t *testing.T) {
			setup(t)

			mockRepository.EXPECT().GetUserBalance(mock.Anything, userID).Return(balance, nil)

			err := service.Withdraw(context.Background(), userID, 100.24, orderNum)

			require.Error(t, err)
			assert.ErrorIs(t, err, domainErrors.ErrOrderNotEnoughBalance)
		},
	)

	t.Run(
		"Error while withdraw", func(t *testing.T) {
			setup(t)

			testError := assert.AnError
			mockRepository.EXPECT().GetUserBalance(mock.Anything, userID).Return(balance, nil)
			mockRepository.EXPECT().Withdraw(
				mock.Anything,
				models.DbBalance{UserID: balance.UserID, ID: balance.ID, Current: 0.01, Withdrawn: 1100.22},
				models.DbWithdrawn{Order: orderNum, Sum: sum, UserID: userID},
			).Return(testError)

			err := service.Withdraw(context.Background(), userID, sum, orderNum)

			require.Error(t, err)
			assert.ErrorIs(t, err, testError)
		},
	)
}
