package middlewares

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run(
		"Successfully", func(t *testing.T) {
			mw, err := New(&zerolog.Logger{})

			require.NoError(t, err)
			require.NotNil(t, mw)
		},
	)

	t.Run(
		"Failed, empty deps", func(t *testing.T) {
			tests := []struct {
				logger *zerolog.Logger
			}{
				{logger: nil},
			}

			for _, test := range tests {
				mw, err := New(test.logger)

				require.Error(t, err)
				require.Nil(t, mw)
			}
		},
	)
}
