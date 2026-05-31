package user

import (
	"context"
	"fmt"
	"time"

	domainErrors "github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Claims структура токена авторизации
type Claims struct {
	UserID int
	Login  string
	jwt.RegisteredClaims
}

var tokenExp = time.Hour

// Auth авторизация пользователя
func (s *Service) Auth(ctx context.Context, login string, password string) (string, error) {
	user, err := s.repository.GetUserByLogin(ctx, login)
	if err != nil {
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", domainErrors.ErrInvalidAuthData
	}

	salt := s.cfg.GetSalt()
	fmt.Println(1)
	jwtString, err := BuildJWTString(*user, salt)
	fmt.Println(2, jwtString)
	fmt.Println(3, err)
	if err != nil {
		return "", err
	}
	return jwtString, nil
}

// BuildJWTString функция получения строкового токена пользователя
func BuildJWTString(user models.DbUser, salt []byte) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
			},
			UserID: user.ID,
			Login:  user.Login,
		},
	)

	tokenString, err := token.SignedString(salt)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
