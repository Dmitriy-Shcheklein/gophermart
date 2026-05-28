package order

import (
	"context"
	"testing"

	domainErrors "github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_Upload(t *testing.T) {
	var (
		mockRepository *MockRepository
		logger         zerolog.Logger
		userID         int
		orderNum       string

		service *Service
	)

	setup := func(t *testing.T) {
		mockRepository = NewMockRepository(t)
		logger = zerolog.Nop()
		userID = 1
		orderNum = "0"

		service, _ = New(&logger, mockRepository, NewMockLoyaltyService(t))
	}

	t.Run(
		"Successfully", func(t *testing.T) {
			setup(t)

			mockRepository.EXPECT().CreateOrder(mock.Anything, userID, orderNum).Return(nil)

			err := service.Upload(context.Background(), userID, orderNum)

			require.NoError(t, err)
		},
	)

	t.Run(
		"Invalid number", func(t *testing.T) {
			setup(t)

			userID = 1
			orderNum = "1234"

			err := service.Upload(context.Background(), userID, orderNum)

			require.Error(t, err)
			assert.ErrorIs(t, err, domainErrors.ErrOrderInvalidNumber)
		},
	)

	t.Run(
		"Order belongs another user", func(t *testing.T) {
			setup(t)

			pgxUniqueError := &pgconn.PgError{Code: "23505"}
			mockRepository.EXPECT().CreateOrder(mock.Anything, userID, orderNum).Return(pgxUniqueError)
			mockRepository.EXPECT().GetOrderByNum(mock.Anything, orderNum).Return(
				models.DbOrder{UserID: userID + 1}, nil,
			)

			err := service.Upload(context.Background(), userID, orderNum)

			require.Error(t, err)
			assert.ErrorIs(t, err, domainErrors.ErrOrderBelongsAnotherUser)
		},
	)

	t.Run(
		"Order already exists", func(t *testing.T) {
			setup(t)

			pgxUniqueError := &pgconn.PgError{Code: "23505"}
			mockRepository.EXPECT().CreateOrder(mock.Anything, userID, orderNum).Return(pgxUniqueError)
			mockRepository.EXPECT().GetOrderByNum(mock.Anything, orderNum).Return(
				models.DbOrder{UserID: userID}, nil,
			)

			err := service.Upload(context.Background(), userID, orderNum)

			require.Error(t, err)
			assert.ErrorIs(t, err, domainErrors.ErrOrderAlreadyExists)
		},
	)

	t.Run(
		"Empty order number error", func(t *testing.T) {
			setup(t)

			mockRepository.EXPECT().CreateOrder(mock.Anything, userID, "").Return(assert.AnError)

			err := service.Upload(context.Background(), userID, "")

			require.Error(t, err)
			assert.Equal(t, assert.AnError, err)
		},
	)

	t.Run(
		"Invalid Luhn algorithm numbers", func(t *testing.T) {
			setup(t)

			invalidNumbers := []string{"123", "12345678901", "123456789012"}

			for _, invalidNum := range invalidNumbers {
				err := service.Upload(context.Background(), userID, invalidNum)

				require.Error(t, err)
				assert.ErrorIs(t, err, domainErrors.ErrOrderInvalidNumber)
			}
		},
	)

	t.Run(
		"Valid Luhn algorithm numbers", func(t *testing.T) {
			setup(t)

			validNumbers := []string{"0", "79927398713"}

			for _, validNum := range validNumbers {
				mockRepository.EXPECT().CreateOrder(mock.Anything, userID, validNum).Return(nil)

				err := service.Upload(context.Background(), userID, validNum)

				require.NoError(t, err)
			}
		},
	)
}
