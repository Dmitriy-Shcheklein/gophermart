package pg

import (
	"context"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
)

const (
	logContext = "userRepository"
)

func (r *Repository) CreateUser(ctx context.Context, login string, password string) (models.DbUser, error) {
	const logCode = "CreateUser"
	r.logger.Debug().Str("context", logContext).Str("code", logCode).Msg("Start")
	query := "insert into users (login, password) values($1, $2) returning id, login, password, created_at"

	var user models.DbUser
	row := r.pool.QueryRow(ctx, query, login, password)
	if err := row.Scan(&user.ID, &user.Login, &user.Password, &user.CreatedAt); err != nil {
		r.logger.Error().Err(err).Str("context", logContext).Str("code", logCode).Msg("Error while register user")
		return models.DbUser{}, err
	}
	r.logger.Debug().Str("context", logContext).Str("code", logCode).Msg("Create user successful")
	return user, nil
}

func (r *Repository) GetUserByLogin(ctx context.Context, login string) (*models.DbUser, error) {
	const logCode = "GetUserByLogin"
	r.logger.Debug().Str("context", logContext).Str("code", logCode).Msg("Start")
	query := "select id, login, password, created_at from users where login = $1"
	row := r.pool.QueryRow(ctx, query, login)
	var result models.DbUser
	err := row.Scan(&result.ID, &result.Login, &result.Password, &result.CreatedAt)
	if err != nil {
		r.logger.Error().Err(err).Str("context", logContext).Str(
			"code", logCode,
		).Msg("Error while getting user by login")
	}
	return &result, err
}
