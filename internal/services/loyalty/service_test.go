package order

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestNewLoyaltyService(t *testing.T) {
	t.Run(
		"Successfully", func(t *testing.T) {
			svc, err := New(&zerolog.Logger{}, NewMockConfig(t), NewMockHttpClient(t))

			require.NoError(t, err)
			require.NotNil(t, svc)
			require.NotNil(t, svc.logger)
			require.NotNil(t, svc.cfg)
			require.NotNil(t, svc.validate)
			require.NotNil(t, svc.httpClient)
		},
	)

	t.Run(
		"Failed, empty deps", func(t *testing.T) {
			logger := &zerolog.Logger{}
			cfg := NewMockConfig(t)
			httpClient := NewMockHttpClient(t)

			tests := []struct {
				logger     *zerolog.Logger
				cfg        Config
				httpClient HttpClient
			}{
				{logger: nil, cfg: cfg, httpClient: httpClient},
				{logger: logger, cfg: nil, httpClient: httpClient},
				{logger: logger, cfg: cfg, httpClient: nil},
			}

			for _, test := range tests {
				svc, err := New(test.logger, test.cfg, test.httpClient)
				require.Error(t, err)
				require.Nil(t, svc)
			}
		},
	)
}
