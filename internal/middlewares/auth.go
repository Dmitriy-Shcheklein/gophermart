package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/services/user"
	"github.com/golang-jwt/jwt/v4"
)

type claims = user.Claims

type token string

const userToken token = "user_token"

var errInvalidUserFormat = errors.New("invalid user format")

// Auth middleware авторизации
func (m *Middleware) Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			jwtToken := r.Header.Get("Authorization")
			if jwtToken == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			withCtx, err := m.verifyToken(r, jwtToken)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			h.ServeHTTP(w, withCtx)
		},
	)
}

func (m *Middleware) verifyToken(r *http.Request, jwtToken string) (*http.Request, error) {
	parts := strings.Split(jwtToken, ".")
	if len(parts) != 3 {
		return r, fmt.Errorf("%w: invalid token format", errInvalidUserFormat)
	}

	currClaims := claims{}
	if _, err := jwt.ParseWithClaims(
		jwtToken, &currClaims, func(t *jwt.Token) (interface{}, error) {
			return m.cfg.GetSalt(), nil
		},
	); err != nil {
		return r, err
	}
	ctx := context.WithValue(r.Context(), userToken, currClaims)
	return r.WithContext(ctx), nil
}

// NewAuthService конструктор для сервиса авторизации
func NewAuthService() *AuthService {
	return &AuthService{}
}

// AuthService структура сервиса авторизации
type AuthService struct{}

// GetUserID метод получения идентификатора пользователя из контекста запроса
func (a *AuthService) GetUserID(ctx context.Context) (int, error) {
	v, ok := ctx.Value(userToken).(claims)
	if !ok {
		return 0, errors.New("error while getting UserID")
	}
	return v.UserID, nil
}
