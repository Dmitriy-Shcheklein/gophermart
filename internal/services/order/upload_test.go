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

		service, _ = New(&logger, mockRepository)
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
}
