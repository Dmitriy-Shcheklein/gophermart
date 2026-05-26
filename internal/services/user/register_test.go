package user

import (
	"context"
	"testing"
	"time"

	domainErrors "github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_Register(t *testing.T) {
	var (
		mockRepository *MockRepository
		logger         *zerolog.Logger
		user           models.DbUser

		service *Service
	)

	setup := func(t *testing.T) {
		mockRepository = NewMockRepository(t)
		nopLogger := zerolog.Nop()
		logger = &nopLogger
		user = models.DbUser{ID: 1, Login: "login", Password: "pass", CreatedAt: time.Now()}

		service, _ = New(logger, mockRepository)
	}

	t.Run(
		"Successfully", func(t *testing.T) {
			setup(t)

			mockRepository.EXPECT().CreateUser(
				mock.Anything, "login", mock.MatchedBy(
					func(hashed string) bool {
						return hashed != "pass"
					},
				),
			).Return(user, nil)

			jwt, err := service.Register(context.Background(), "login", "pass")

			require.NoError(t, err)
			assert.NotEmpty(t, jwt)
		},
	)

	t.Run(
		"No rows error", func(t *testing.T) {
			setup(t)

			mockRepository.EXPECT().CreateUser(mock.Anything, mock.Anything, mock.Anything).Return(user, pgx.ErrNoRows)

			_, err := service.Register(context.Background(), "login", "pass")

			require.Error(t, err)
			assert.Equal(t, domainErrors.ErrLoginDuplicate, err)
		},
	)

	t.Run(
		"Unexpected error", func(t *testing.T) {
			setup(t)

			testError := assert.AnError
			mockRepository.EXPECT().CreateUser(mock.Anything, mock.Anything, mock.Anything).Return(user, testError)

			_, err := service.Register(context.Background(), "login", "pass")

			require.Error(t, err)
			assert.Equal(t, testError, err)
		},
	)
}
