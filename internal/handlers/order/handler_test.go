package user

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUserHandler(t *testing.T) {
	t.Run(
		"Successfully", func(t *testing.T) {
			handler, err := New(&zerolog.Logger{}, NewMockService(t), NewMockAuthService(t))

			require.NoError(t, err)
			require.NotNil(t, handler)
			assert.NotNil(t, handler.validate)
			assert.NotNil(t, handler.authService)
			assert.NotNil(t, handler.service)
		},
	)

	t.Run(
		"Failed, empty deps", func(t *testing.T) {
			logger := &zerolog.Logger{}
			service := NewMockService(t)
			authSvc := NewMockAuthService(t)

			tests := []struct {
				logger  *zerolog.Logger
				service Service
				authSvc AuthService
			}{
				{logger: nil, service: service, authSvc: authSvc},
				{logger: logger, service: nil, authSvc: authSvc},
				{logger: logger, service: service, authSvc: nil},
			}

			for _, test := range tests {
				handler, err := New(test.logger, test.service, test.authSvc)
				require.Error(t, err)
				require.Nil(t, handler)
			}
		},
	)
}
