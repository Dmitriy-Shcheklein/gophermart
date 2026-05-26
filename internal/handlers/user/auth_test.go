package user

import (
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

func TestHandler_Auth(t *testing.T) {
	var (
		mockService *MockService
		w           *httptest.ResponseRecorder
		r           *http.Request
		body        *strings.Reader

		handler *Handler
	)

	setup := func(t *testing.T) {
		logger := zerolog.Nop()
		mockService = NewMockService(t)
		body = strings.NewReader("{\"login\": \"firstLogin\",\"password\": \"pass\"}")
		r = httptest.NewRequest(http.MethodPost, "/", body)
		r.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		handler, _ = New(&logger, mockService)
	}

	t.Run(
		"Successfully", func(t *testing.T) {
			setup(t)

			token := "token"
			mockService.EXPECT().Auth(mock.Anything, "firstLogin", "pass").Return(token, nil)

			handler.Auth(w, r)

			require.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, token, w.Header().Get("Authorization"))
		},
	)

	t.Run(
		"Invalid content-type", func(t *testing.T) {
			setup(t)
			r.Header = http.Header{}

			handler.Auth(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			assert.Equal(t, domainErrors.InvalidContentTypeMsg, strings.TrimSpace(w.Body.String()))
		},
	)

	t.Run(
		"Decode body error", func(t *testing.T) {
			setup(t)

			body = strings.NewReader("{\"login\": \"firstLogin,\"password\": \"pass\"}")
			r = httptest.NewRequest(http.MethodPost, "/", body)
			r.Header.Set("Content-Type", "application/json")

			handler.Auth(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			assert.Equal(t, domainErrors.DecodeBodyErrMsg, strings.TrimSpace(w.Body.String()))
		},
	)

	t.Run(
		"Validate body error", func(t *testing.T) {
			setup(t)

			body = strings.NewReader("{\"login\": \"\",\"password\": \"pass\"}")
			r = httptest.NewRequest(http.MethodPost, "/", body)
			r.Header.Set("Content-Type", "application/json")

			handler.Auth(w, r)

			require.Equal(t, http.StatusBadRequest, w.Code)
			assert.Equal(t, domainErrors.ValidateBodyErrMsg, strings.TrimSpace(w.Body.String()))
		},
	)

	t.Run(
		"Invalid auth data error", func(t *testing.T) {
			setup(t)

			mockService.EXPECT().Auth(
				mock.Anything, mock.Anything, mock.Anything,
			).Return("", domainErrors.ErrInvalidAuthData)

			handler.Auth(w, r)

			require.Equal(t, http.StatusUnauthorized, w.Code)
		},
	)

	t.Run(
		"Internal error", func(t *testing.T) {
			setup(t)

			mockService.EXPECT().Auth(mock.Anything, mock.Anything, mock.Anything).Return("", assert.AnError)

			handler.Auth(w, r)

			require.Equal(t, http.StatusInternalServerError, w.Code)
		},
	)
}
