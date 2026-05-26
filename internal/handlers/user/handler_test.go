package user

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestNewUserHandler(t *testing.T) {
	t.Run(
		"Successfully", func(t *testing.T) {
			handler, err := New(&zerolog.Logger{}, NewMockService(t))

			require.NoError(t, err)
			require.NotNil(t, handler)
			require.NotNil(t, handler.validate)
		},
	)

	t.Run(
		"Failed, empty deps", func(t *testing.T) {
			logger := &zerolog.Logger{}
			service := NewMockService(t)

			tests := []struct {
				logger  *zerolog.Logger
				service Service
			}{
				{logger: nil, service: service},
				{logger: logger, service: nil},
			}

			for _, test := range tests {
				handler, err := New(test.logger, nil)
				require.Error(t, err)
				require.Nil(t, handler)
			}
		},
	)
}
