package order

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_GetBalance(t *testing.T) {
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

			svcResult := models.ResponseBalance{
				Current:   10,
				Withdrawn: 100,
			}

			mockAuthService.EXPECT().GetUserID(ctx).Return(userID, nil)
			mockService.EXPECT().GetBalance(ctx, userID).Return(svcResult, nil)

			handler.GetBalance(w, r)

			require.Equal(t, http.StatusOK, w.Code)
			assert.Equal(
				t,
				"{\"current\":10,\"withdrawn\":100}",
				strings.TrimSpace(w.Body.String()),
			)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		},
	)

	t.Run(
		"Error while get userID", func(t *testing.T) {
			setup(t)

			mockAuthService.EXPECT().GetUserID(ctx).Return(0, assert.AnError)

			handler.GetBalance(w, r)

			require.Equal(t, http.StatusInternalServerError, w.Code)
		},
	)

	t.Run(
		"Error from GetBalance", func(t *testing.T) {
			setup(t)

			mockAuthService.EXPECT().GetUserID(ctx).Return(userID, nil)
			mockService.EXPECT().GetBalance(ctx, userID).Return(models.ResponseBalance{}, assert.AnError)

			handler.GetBalance(w, r)

			require.Equal(t, http.StatusInternalServerError, w.Code)
		},
	)
}
