package user

import (
	"context"
	"testing"
	"time"

	domainErrors "github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_Auth(t *testing.T) {
	var (
		mockRepository *MockRepository
		logger         *zerolog.Logger
		user           models.DbUser
		password       string

		service *Service
	)

	setup := func(t *testing.T) {
		mockRepository = NewMockRepository(t)
		nopLogger := zerolog.Nop()
		logger = &nopLogger

		password, _ = hashPassword("pass")
		user = models.DbUser{ID: 1, Login: "login", Password: password, CreatedAt: time.Now()}

		service, _ = New(logger, mockRepository)
	}

	t.Run(
		"Successfully", func(t *testing.T) {
			setup(t)

			mockRepository.EXPECT().GetUserByLogin(mock.Anything, "login").Return(&user, nil)

			jwt, err := service.Auth(context.Background(), "login", "pass")

			require.NoError(t, err)
			require.NotEmpty(t, jwt)
		},
	)

	t.Run(
		"Repository error", func(t *testing.T) {
			setup(t)

			testError := assert.AnError
			mockRepository.EXPECT().GetUserByLogin(mock.Anything, mock.Anything).Return(nil, testError)

			_, err := service.Auth(context.Background(), "login", "pass")

			require.Error(t, err)
			assert.Equal(t, testError, err)
		},
	)

	t.Run(
		"Invalid login/pass error", func(t *testing.T) {
			setup(t)

			user.Password = "random"
			mockRepository.EXPECT().GetUserByLogin(mock.Anything, mock.Anything).Return(&user, nil)

			_, err := service.Auth(context.Background(), "login", "pass")

			require.Error(t, err)
			assert.Equal(t, domainErrors.ErrInvalidAuthData, err)
		},
	)
}
