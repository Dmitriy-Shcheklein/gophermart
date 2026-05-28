package user

import (
	"context"
	"testing"
	"time"

	domainErrors "github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
				), models.DbBalance{Withdrawn: 0, Current: 0},
			).Return(user, nil)

			jwt, err := service.Register(context.Background(), "login", "pass")

			require.NoError(t, err)
			assert.NotEmpty(t, jwt)
		},
	)

	t.Run(
		"No rows error", func(t *testing.T) {
			setup(t)

			mockRepository.EXPECT().CreateUser(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
				user, pgx.ErrNoRows,
			)

			_, err := service.Register(context.Background(), "login", "pass")

			require.Error(t, err)
			assert.Equal(t, domainErrors.ErrLoginDuplicate, err)
		},
	)

	t.Run(
		"Unexpected error", func(t *testing.T) {
			setup(t)

			testError := assert.AnError
			mockRepository.EXPECT().CreateUser(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
				user, testError,
			)

			_, err := service.Register(context.Background(), "login", "pass")

			require.Error(t, err)
			assert.Equal(t, testError, err)
		},
	)

	t.Run(
		"Duplicate user error", func(t *testing.T) {
			setup(t)

			mockRepository.EXPECT().CreateUser(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(user, &pgconn.PgError{Code: "23505"})

			_, err := service.Register(context.Background(), "existinguser", "password")

			require.Error(t, err)
		},
	)

	t.Run(
		"Database error on user creation", func(t *testing.T) {
			setup(t)

			mockRepository.EXPECT().CreateUser(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(user, assert.AnError)

			_, err := service.Register(context.Background(), "testuser", "password")

			require.Error(t, err)
			assert.Equal(t, assert.AnError, err)
		},
	)

	t.Run(
		"Password hashing error", func(t *testing.T) {
			setup(t)

			mockRepository.EXPECT().CreateUser(mock.Anything, "", mock.Anything, mock.Anything).Return(
				user, pgx.ErrNoRows,
			)

			_, err := service.Register(context.Background(), "", "password")

			require.Error(t, err)
			assert.Equal(t, domainErrors.ErrLoginDuplicate, err)
		},
	)
}
