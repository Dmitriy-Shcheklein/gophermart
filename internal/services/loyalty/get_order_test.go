package order

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	domainErrors "github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_GetOrder(t *testing.T) {
	var (
		logger         *zerolog.Logger
		mockCfg        *MockConfig
		mockHttpClient *MockHttpClient
		body           io.ReadCloser

		service *Service
	)

	setup := func(t *testing.T) {
		nopLogger := zerolog.Nop()
		logger = &nopLogger
		mockCfg = NewMockConfig(t)
		mockHttpClient = NewMockHttpClient(t)
		bodyData := "{\"order\": \"1\",\"status\": \"PROCESSED\",\"accrual\": 100}"
		body = io.NopCloser(bytes.NewReader([]byte(bodyData)))
		service, _ = New(logger, mockCfg, mockHttpClient)
	}

	t.Run(
		"Successfully", func(t *testing.T) {
			setup(t)

			mockCfg.EXPECT().GetAccrualSrvAddr().Return("address")
			mockHttpClient.EXPECT().Get("address/api/orders/1").Return(
				&http.Response{StatusCode: http.StatusOK, Body: body}, nil,
			)

			result, err := service.GetOrder(context.Background(), "1")

			require.NoError(t, err)
			assert.Equal(
				t, "1", result.Order,
			)
			assert.Equal(
				t, models.ProcessedOrder, result.Status,
			)
			assert.EqualValues(
				t, 100, *result.Accrual,
			)
		},
	)

	t.Run(
		"Is Wait", func(t *testing.T) {
			setup(t)

			service.isWait.Store(true)

			_, err := service.GetOrder(context.Background(), "1")

			require.ErrorIs(t, err, domainErrors.ErrorLoyaltyWait)
		},
	)

	t.Run(
		"http client error", func(t *testing.T) {
			setup(t)

			testError := assert.AnError
			mockCfg.EXPECT().GetAccrualSrvAddr().Return("address")
			mockHttpClient.EXPECT().Get("address/api/orders/1").Return(
				nil, testError,
			)

			_, err := service.GetOrder(context.Background(), "1")

			require.ErrorIs(t, err, testError)
		},
	)

	t.Run(
		"Staus No-Content", func(t *testing.T) {
			setup(t)

			mockCfg.EXPECT().GetAccrualSrvAddr().Return("address")
			mockHttpClient.EXPECT().Get("address/api/orders/1").Return(
				&http.Response{StatusCode: http.StatusNoContent}, nil,
			)

			result, err := service.GetOrder(context.Background(), "1")

			require.NoError(t, err)
			assert.Equal(
				t, "1", result.Order,
			)
			assert.Equal(
				t, models.InvalidOrder, result.Status,
			)
			assert.Nil(
				t, result.Accrual,
			)
		},
	)

	t.Run(
		"Status Internal Server Error", func(t *testing.T) {
			setup(t)

			mockCfg.EXPECT().GetAccrualSrvAddr().Return("address")
			mockHttpClient.EXPECT().Get("address/api/orders/1").Return(
				&http.Response{StatusCode: http.StatusInternalServerError}, nil,
			)

			result, err := service.GetOrder(context.Background(), "1")

			require.ErrorIs(t, err, domainErrors.ErrLoyaltyUnknown)
			assert.Empty(t, result)
		},
	)

	t.Run(
		"Unknown status", func(t *testing.T) {
			setup(t)

			mockCfg.EXPECT().GetAccrualSrvAddr().Return("address")
			mockHttpClient.EXPECT().Get("address/api/orders/1").Return(
				&http.Response{StatusCode: http.StatusAlreadyReported}, nil,
			)

			result, err := service.GetOrder(context.Background(), "1")

			require.ErrorIs(t, err, domainErrors.ErrLoyaltyUnknownStatusCode)
			assert.Empty(t, result)
		},
	)

	t.Run(
		"Validation error", func(t *testing.T) {
			setup(t)

			body = io.NopCloser(bytes.NewReader([]byte("{\"order\": \"1\",\"status\": \"START\",\"accrual\": 100}")))

			mockCfg.EXPECT().GetAccrualSrvAddr().Return("address")
			mockHttpClient.EXPECT().Get("address/api/orders/1").Return(
				&http.Response{StatusCode: http.StatusOK, Body: body}, nil,
			)

			result, err := service.GetOrder(context.Background(), "1")

			require.ErrorIs(t, err, domainErrors.ErrLoyaltyValidateBody)
			assert.Empty(t, result)
		},
	)

	t.Run(
		"Decode error", func(t *testing.T) {
			setup(t)

			body = io.NopCloser(bytes.NewReader([]byte("random")))

			mockCfg.EXPECT().GetAccrualSrvAddr().Return("address")
			mockHttpClient.EXPECT().Get("address/api/orders/1").Return(
				&http.Response{StatusCode: http.StatusOK, Body: body}, nil,
			)

			result, err := service.GetOrder(context.Background(), "1")

			require.ErrorIs(t, err, domainErrors.ErrLoyaltyDecodeBody)
			assert.Empty(t, result)
		},
	)

	t.Run(
		"Handle 429 status code", func(t *testing.T) {
			setup(t)

			withCancel, cancelFunc := context.WithCancel(context.Background())
			defer cancelFunc()

			response := &http.Response{StatusCode: http.StatusTooManyRequests}
			response.Header = http.Header{}
			response.Header.Set("Retry-After", "10")

			mockCfg.EXPECT().GetAccrualSrvAddr().Return("address")
			mockHttpClient.EXPECT().Get("address/api/orders/1").Return(
				response, nil,
			)

			result, err := service.GetOrder(withCancel, "1")

			require.ErrorIs(t, err, domainErrors.ErrorLoyaltyTooManyRequest)
			assert.Empty(t, result)
			assert.True(t, service.isWait.Load())
		},
	)
}
