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
	"github.com/stretchr/testify/require"
)

func TestHandler_Upload(t *testing.T) {
	var (
		mockService     *MockService
		mockAuthService *MockAuthService
		w               *httptest.ResponseRecorder
		r               *http.Request
		body            *strings.Reader
		orderNum        string
		userID          int
		ctx             context.Context

		handler *Handler
	)

	setup := func(t *testing.T) {
		logger := zerolog.Nop()
		mockService = NewMockService(t)
		mockAuthService = NewMockAuthService(t)
		orderNum = "123456"
		userID = 1
		ctx = context.Background()
		body = strings.NewReader(orderNum)
		r = httptest.NewRequestWithContext(ctx, http.MethodPost, "/", body)
		r.Header.Set("Content-Type", "text/plain")
		w = httptest.NewRecorder()

		handler, _ = New(&logger, mockService, mockAuthService)
	}

	t.Run(
		"Successfully", func(t *testing.T) {
			setup(t)

			mockAuthService.EXPECT().GetUserID(ctx).Return(userID, nil)
			mockService.EXPECT().Upload(ctx, userID, orderNum).Return(nil)

			handler.Upload(w, r)

			require.Equal(t, http.StatusAccepted, w.Code)
		},
	)

	t.Run(
		"Already exists", func(t *testing.T) {
			setup(t)

			mockAuthService.EXPECT().GetUserID(ctx).Return(userID, nil)
			mockService.EXPECT().Upload(ctx, userID, orderNum).Return(domainErrors.ErrOrderAlreadyExists)

			handler.Upload(w, r)

			require.Equal(t, http.StatusOK, w.Code)
		},
	)

	t.Run(
		"Belongs another user", func(t *testing.T) {
			setup(t)

			mockAuthService.EXPECT().GetUserID(ctx).Return(userID, nil)
			mockService.EXPECT().Upload(ctx, userID, orderNum).Return(domainErrors.ErrOrderBelongsAnotherUser)

			handler.Upload(w, r)

			require.Equal(t, http.StatusConflict, w.Code)
		},
	)

	t.Run(
		"Invalid number", func(t *testing.T) {
			setup(t)

			mockAuthService.EXPECT().GetUserID(ctx).Return(userID, nil)
			mockService.EXPECT().Upload(ctx, userID, orderNum).Return(domainErrors.ErrOrderInvalidNumber)

			handler.Upload(w, r)

			require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		},
	)

	t.Run(
		"Unexpected error", func(t *testing.T) {
			setup(t)

			mockAuthService.EXPECT().GetUserID(ctx).Return(userID, nil)
			mockService.EXPECT().Upload(ctx, userID, orderNum).Return(assert.AnError)

			handler.Upload(w, r)

			require.Equal(t, http.StatusInternalServerError, w.Code)
		},
	)

	t.Run(
		"error while getting userID", func(t *testing.T) {
			setup(t)

			mockAuthService.EXPECT().GetUserID(ctx).Return(0, assert.AnError)

			handler.Upload(w, r)

			require.Equal(t, http.StatusInternalServerError, w.Code)
		},
	)

	t.Run(
		"invalid content type", func(t *testing.T) {
			setup(t)

			r = httptest.NewRequestWithContext(ctx, http.MethodPost, "/", body)
			r.Header.Set("Content-Type", "application/json")

			handler.Upload(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), domainErrors.InvalidContentTypeMsg)
		},
	)

	t.Run(
		"invalid validate order num", func(t *testing.T) {
			setup(t)

			body = strings.NewReader("")
			r = httptest.NewRequestWithContext(ctx, http.MethodPost, "/", body)
			r.Header.Set("Content-Type", "text/plain")

			handler.Upload(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), domainErrors.ValidateBodyErrMsg)
		},
	)
}
