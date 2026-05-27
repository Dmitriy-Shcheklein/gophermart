package order

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_GetByUserID(t *testing.T) {
	var (
		mockRepository *MockRepository
		logger         zerolog.Logger
		userID         int
		repoResult     []models.DbOrder
		uploadTime     time.Time
		firstAccrual   float64

		service *Service
	)

	setup := func(t *testing.T) {
		mockRepository = NewMockRepository(t)
		logger = zerolog.Nop()
		userID = 1
		uploadTime = time.Now()
		firstAccrual = 100.2

		repoResult = []models.DbOrder{
			{
				ID:         1,
				Number:     "123",
				Accrual:    sql.NullFloat64{Float64: firstAccrual, Valid: true},
				UserID:     userID,
				UploadedAt: uploadTime,
				Status:     models.NewOrder,
			},
			{
				ID:         2,
				Number:     "321",
				Accrual:    sql.NullFloat64{},
				UserID:     userID,
				UploadedAt: uploadTime,
				Status:     models.ProcessedOrder,
			},
		}

		service, _ = New(&logger, mockRepository, NewMockLoyaltyService(t))
	}

	t.Run(
		"Successfully", func(t *testing.T) {
			setup(t)

			mockRepository.EXPECT().GetByUserId(mock.Anything, userID).Return(repoResult, nil)

			res, err := service.GetList(context.Background(), userID)

			require.NoError(t, err)
			assert.Equal(
				t, []models.RequestOrder{
					{
						Number:     "123",
						Accrual:    &firstAccrual,
						UploadedAt: uploadTime,
						Status:     models.NewOrder,
					},
					{
						Number:     "321",
						Accrual:    nil,
						UploadedAt: uploadTime,
						Status:     models.ProcessedOrder,
					},
				}, res,
			)
		},
	)

	t.Run(
		"Repository err", func(t *testing.T) {
			setup(t)

			testError := assert.AnError
			mockRepository.EXPECT().GetByUserId(mock.Anything, userID).Return(nil, testError)

			_, err := service.GetList(context.Background(), userID)

			require.Error(t, err)
			assert.Equal(t, testError, err)
		},
	)
}
