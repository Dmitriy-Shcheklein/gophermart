package user

import (
	"context"
	"errors"

	domainErrors "github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/jackc/pgx/v5"
)

func (s *Service) Register(ctx context.Context, login, password string) error {
	hashPass, err := hashPassword(password)
	if err != nil {
		return err
	}
	if err = s.repository.CreateUser(ctx, login, hashPass); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainErrors.ErrLoginDuplicate
		}
		return err
	}
	return nil
}
