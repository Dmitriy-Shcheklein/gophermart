package user

import (
	"context"
	"errors"

	domainErrors "github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
	"github.com/jackc/pgx/v5"
)

func (s *Service) Register(ctx context.Context, login, password string) (string, error) {
	hashPass, err := hashPassword(password)
	if err != nil {
		return "", err
	}
	dbUSer, err := s.repository.CreateUser(ctx, login, hashPass, models.DbBalance{Withdrawn: 0, Current: 0})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", domainErrors.ErrLoginDuplicate
		}
		return "", err
	}
	return BuildJWTString(dbUSer)
}
