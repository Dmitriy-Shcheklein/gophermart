package pg

import (
	"context"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
)

const (
	logContext = "userRepository"
)

func (r *Repository) CreateUser(ctx context.Context, login string, password string) error {
	const logCode = "CreateUser"
	r.logger.Debug().Str("context", logContext).Str("code", logCode).Msg("Start")
	query := "insert into users (login, password) values($1, $2)"

	_, err := r.pool.Exec(ctx, query, login, password)
	if err != nil {
		r.logger.Error().Err(err).Str("context", logContext).Str("code", logCode).Msg("Error while register user")
		return err
	}
	r.logger.Debug().Str("context", logContext).Str("code", logCode).Msg("Create user successful")
	return nil
}

func (r *Repository) GetUserByLogin(ctx context.Context, login string) (*models.DbUser, error) {
	const logCode = "GetUserByLogin"
	r.logger.Debug().Str("context", logContext).Str("code", logCode).Msg("Start")
	query := "select id, login, password, crfeated_at from users where login = $1"
	row := r.pool.QueryRow(ctx, query, login)
	var result models.DbUser
	err := row.Scan(&result)
	if err != nil {
		r.logger.Error().Err(err).Str("context", logContext).Str(
			"code", logCode,
		).Msg("Error while getting user by login")
	}
	return &result, err
}
