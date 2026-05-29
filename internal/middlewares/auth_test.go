package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
	userSvc "github.com/Dmitriy-Shcheklein/gophermart/internal/services/user"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestMiddleware_Auth(t *testing.T) {
	var (
		logger      zerolog.Logger
		nextHandler http.Handler
		r           *http.Request
		w           *httptest.ResponseRecorder
		user        models.DbUser
		jwtString   string
		mockCfg     *MockConfig
		salt        []byte

		mw *Middleware
	)

	setup := func(t *testing.T) {
		logger = zerolog.Nop()
		w = httptest.NewRecorder()
		nextHandler = http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("response"))
			},
		)
		user = models.DbUser{
			ID:        1,
			Login:     "login",
			Password:  "pass",
			CreatedAt: time.Now(),
		}
		salt = []byte("salt")
		jwtStr, err := userSvc.BuildJWTString(user, salt)
		require.NoError(t, err)
		jwtString = jwtStr
		mockCfg = NewMockConfig(t)

		mw, _ = New(&logger, mockCfg)
	}

	t.Run(
		"Successfully", func(t *testing.T) {
			setup(t)

			r = httptest.NewRequest(http.MethodPost, "/", nil)
			r.Header.Set("Authorization", jwtString)

			mockCfg.EXPECT().GetSalt().Return(salt)

			handler := mw.Auth(nextHandler)
			handler.ServeHTTP(w, r)

			require.Equal(t, http.StatusOK, w.Code)
		},
	)

	t.Run(
		"Header is not exists", func(t *testing.T) {
			setup(t)

			r = httptest.NewRequest(http.MethodPost, "/", nil)

			handler := mw.Auth(nextHandler)
			handler.ServeHTTP(w, r)

			require.Equal(t, http.StatusUnauthorized, w.Code)
		},
	)

	t.Run(
		"Header is empty", func(t *testing.T) {
			setup(t)

			r = httptest.NewRequest(http.MethodPost, "/", nil)
			r.Header.Set("Authorization", "")

			handler := mw.Auth(nextHandler)
			handler.ServeHTTP(w, r)

			require.Equal(t, http.StatusUnauthorized, w.Code)
		},
	)

	t.Run(
		"Invalid token format", func(t *testing.T) {
			setup(t)

			r = httptest.NewRequest(http.MethodPost, "/", nil)
			r.Header.Set("Authorization", "first.second.third.fourth")

			handler := mw.Auth(nextHandler)
			handler.ServeHTTP(w, r)

			require.Equal(t, http.StatusUnauthorized, w.Code)
			require.Contains(t, w.Body.String(), "invalid token format")
		},
	)

	t.Run(
		"Invalid token format", func(t *testing.T) {
			setup(t)

			r = httptest.NewRequest(http.MethodPost, "/", nil)
			r.Header.Set("Authorization", "first.second")

			handler := mw.Auth(nextHandler)
			handler.ServeHTTP(w, r)

			require.Equal(t, http.StatusUnauthorized, w.Code)
			require.Contains(t, w.Body.String(), "invalid token format")
		},
	)
}
