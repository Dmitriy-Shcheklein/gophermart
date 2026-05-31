package middlewares

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run(
		"Successfully", func(t *testing.T) {
			mw, err := New(&zerolog.Logger{}, NewMockConfig(t))

			require.NoError(t, err)
			require.NotNil(t, mw)
			require.NotNil(t, mw.logger)
			require.NotNil(t, mw.cfg)
		},
	)

	t.Run(
		"Failed, empty deps", func(t *testing.T) {
			logger := &zerolog.Logger{}
			cfg := NewMockConfig(t)

			tests := []struct {
				logger *zerolog.Logger
				cfg    Config
			}{
				{logger: nil, cfg: cfg},
				{logger: logger, cfg: nil},
			}

			for _, test := range tests {
				mw, err := New(test.logger, test.cfg)

				require.Error(t, err)
				require.Nil(t, mw)
			}
		},
	)
}
