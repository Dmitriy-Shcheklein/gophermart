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

var secretKey = user.SecretKey

type Claims = user.Claims

type UserToken string

const userToken UserToken = "user_token"

var errInvalidUserFormat = errors.New("invalid user format")

func (m *Middleware) Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			jwtToken := r.Header.Get("Authorization")
			if jwtToken == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			withCtx, err := verifyToken(r, jwtToken)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			h.ServeHTTP(w, withCtx)
		},
	)
}

func verifyToken(r *http.Request, jwtToken string) (*http.Request, error) {
	parts := strings.Split(jwtToken, ".")
	if len(parts) != 3 {
		return r, fmt.Errorf("%w: invalid token format", errInvalidUserFormat)
	}

	claims := &Claims{}
	if _, err := jwt.ParseWithClaims(
		jwtToken, claims, func(t *jwt.Token) (interface{}, error) {
			return secretKey, nil
		},
	); err != nil {
		return r, err
	}
	ctx := context.WithValue(r.Context(), userToken, claims)
	return r.WithContext(ctx), nil
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

type AuthService struct{}

func (a *AuthService) GetUserID(ctx context.Context) (int, error) {
	v, ok := ctx.Value(userToken).(Claims)
	if !ok {
		return 0, errors.New("error while getting UserID")
	}
	return v.UserID, nil
}
