package worker

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run(
		"Successfully", func(t *testing.T) {
			svc, err := New(&zerolog.Logger{}, NewMockService(t))

			require.NoError(t, err)
			require.NotNil(t, svc)
			require.NotNil(t, svc.logger)
			require.NotNil(t, svc.service)
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
				svc, err := New(test.logger, test.service)
				require.Error(t, err)
				require.Nil(t, svc)
			}
		},
	)
}

func TestWorker(t *testing.T) {
	var (
		mockService *MockService
		logger      *zerolog.Logger

		worker *Worker
	)

	setup := func(t *testing.T) {
		nopLogger := zerolog.Nop()
		logger = &nopLogger
		mockService = NewMockService(t)

		worker, _ = New(logger, mockService)
	}

	t.Run(
		"Start successfully", func(t *testing.T) {
			setup(t)
			ctx, cancelFunc := context.WithCancel(context.Background())

			mockService.EXPECT().ProcessOrders(mock.Anything).RunAndReturn(
				func(_ context.Context) error {
					cancelFunc()
					return nil
				},
			)
			worker.Start(ctx, time.Millisecond)

			select {
			case <-ctx.Done():
				assert.Equal(t, time.Millisecond, worker.timeout)
				assert.True(t, worker.isStarted.Load())
				require.NotNil(t, worker.cancelFunc)
			case <-time.After(1 * time.Second):
				t.Fatal("Превышен таймаут")
			}
		},
	)

	t.Run(
		"Stop Successfully", func(t *testing.T) {
			setup(t)

			ctx := context.Background()

			mockService.EXPECT().ProcessOrders(mock.Anything).Return(nil)

			worker.Start(ctx, time.Hour)

			time.Sleep(time.Millisecond * 100)

			worker.Stop()

			require.False(t, worker.isStarted.Load())
		},
	)
}
