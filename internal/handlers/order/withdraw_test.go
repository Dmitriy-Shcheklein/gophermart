package order

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	domainErrors "github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_Withdraw(t *testing.T) {
	var (
		mockService     *MockService
		mockAuthService *MockAuthService
		w               *httptest.ResponseRecorder
		r               *http.Request
		body            *strings.Reader
		userID          int
		ctx             context.Context

		handler *Handler
	)

	setup := func(t *testing.T) {
		logger := zerolog.Nop()
		mockService = NewMockService(t)
		mockAuthService = NewMockAuthService(t)
		userID = 1
		ctx = context.Background()
		body = strings.NewReader("{\"sum\": 100.22,\"order\": \"1234560\"}")
		r = httptest.NewRequestWithContext(ctx, http.MethodPost, "/", body)
		r.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		handler, _ = New(&logger, mockService, mockAuthService)
	}

	t.Run(
		"Successfully", func(t *testing.T) {
			setup(t)

			mockAuthService.EXPECT().GetUserID(ctx).Return(userID, nil)
			mockService.EXPECT().Withdraw(ctx, userID, 100.22, "1234560").Return(nil)

			handler.Withdraw(w, r)

			require.Equal(t, http.StatusOK, w.Code)
		},
	)

	t.Run(
		"Invalid content/type", func(t *testing.T) {
			setup(t)

			r = httptest.NewRequestWithContext(ctx, http.MethodPost, "/", body)
			r.Header.Set("Content-Type", "text/plain")

			handler.Withdraw(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), domainErrors.InvalidContentTypeMsg)
		},
	)

	t.Run(
		"Error while getting userID", func(t *testing.T) {
			setup(t)

			mockAuthService.EXPECT().GetUserID(ctx).Return(userID, assert.AnError)

			handler.Withdraw(w, r)

			require.Equal(t, http.StatusInternalServerError, w.Code)
		},
	)

	t.Run(
		"Error while call svc Withdraw", func(t *testing.T) {
			setup(t)

			mockAuthService.EXPECT().GetUserID(ctx).Return(userID, nil)
			mockService.EXPECT().Withdraw(
				mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			).Return(assert.AnError)

			handler.Withdraw(w, r)

			require.Equal(t, http.StatusInternalServerError, w.Code)
		},
	)

	t.Run(
		"Not enough balance", func(t *testing.T) {
			setup(t)

			mockAuthService.EXPECT().GetUserID(ctx).Return(userID, nil)
			mockService.EXPECT().Withdraw(
				mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			).Return(domainErrors.ErrOrderNotEnoughBalance)

			handler.Withdraw(w, r)

			require.Equal(t, http.StatusPaymentRequired, w.Code)
		},
	)

	t.Run(
		"Order invalid number", func(t *testing.T) {
			setup(t)

			mockAuthService.EXPECT().GetUserID(ctx).Return(userID, nil)
			mockService.EXPECT().Withdraw(
				mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			).Return(domainErrors.ErrOrderInvalidNumber)

			handler.Withdraw(w, r)

			require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		},
	)
}
