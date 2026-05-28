package user

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_GetWithdrawals(t *testing.T) {
	var (
		mockService     *MockService
		mockAuthService *MockAuthService
		w               *httptest.ResponseRecorder
		r               *http.Request
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
		r = httptest.NewRequestWithContext(ctx, http.MethodPost, "/", nil)
		w = httptest.NewRecorder()

		handler, _ = New(&logger, mockService, mockAuthService)
	}

	t.Run(
		"Successfully", func(t *testing.T) {
			setup(t)

			processedAt, err := time.Parse(time.RFC3339Nano, "2026-05-26T22:17:19.31+05:00")
			if err != nil {
				require.NoError(t, err)
			}

			svcResult := []models.ResponseWithdrawn{
				{
					Sum:         1,
					Order:       "123",
					ProcessedAt: processedAt,
				},
			}

			mockAuthService.EXPECT().GetUserID(ctx).Return(userID, nil)
			mockService.EXPECT().GetWithdrawals(ctx, userID).Return(svcResult, nil)

			handler.GetWithdrawals(w, r)

			require.Equal(t, http.StatusOK, w.Code)
			assert.Equal(
				t,
				"[{\"sum\":1,\"order\":\"123\",\"processed_at\":\"2026-05-26T22:17:19.31+05:00\"}]",
				strings.TrimSpace(w.Body.String()),
			)
		},
	)

	t.Run(
		"Empty list", func(t *testing.T) {
			setup(t)

			svcResult := make([]models.ResponseWithdrawn, 0)

			mockAuthService.EXPECT().GetUserID(ctx).Return(userID, nil)
			mockService.EXPECT().GetWithdrawals(ctx, userID).Return(svcResult, nil)

			handler.GetWithdrawals(w, r)

			require.Equal(t, http.StatusNoContent, w.Code)
		},
	)

	t.Run(
		"Error while get userID", func(t *testing.T) {
			setup(t)

			mockAuthService.EXPECT().GetUserID(ctx).Return(0, assert.AnError)

			handler.GetWithdrawals(w, r)

			require.Equal(t, http.StatusInternalServerError, w.Code)
		},
	)

	t.Run(
		"Error from GetList", func(t *testing.T) {
			setup(t)

			mockAuthService.EXPECT().GetUserID(ctx).Return(userID, nil)
			mockService.EXPECT().GetWithdrawals(ctx, userID).Return(nil, assert.AnError)

			handler.GetWithdrawals(w, r)

			require.Equal(t, http.StatusInternalServerError, w.Code)
		},
	)
}
