package postgres

import (
	"context"
	"fmt"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
)

const (
	logContext = "userRepository"
)

// CreateUser метод создания нового пользователя
func (r *Repository) CreateUser(
	ctx context.Context, login string, password string, balance models.DbBalance,
) (models.DbUser, error) {
	const logCode = "CreateUser"
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return models.DbUser{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if err = tx.Rollback(ctx); err != nil {
				r.logger.Error().Err(err).Msg("error while rollback transaction")
			}
		}
	}()
	r.logger.Debug().Str("context", logContext).Str("code", logCode).Msg("Start")
	query := "insert into users (login, password) values($1, $2) returning id, login, password, created_at"

	var user models.DbUser
	row := tx.QueryRow(ctx, query, login, password)
	if err = row.Scan(&user.ID, &user.Login, &user.Password, &user.CreatedAt); err != nil {
		r.logger.Error().Err(err).Str("context", logContext).Str("code", logCode).Msg("Error while register user")
		return models.DbUser{}, err
	}
	balanceQuery := "insert into balances (current, withdrawn, user_id) values ($1, $2, $3)"
	_, err = tx.Exec(ctx, balanceQuery, balance.Current, balance.Withdrawn, user.ID)
	if err != nil {
		return models.DbUser{}, err
	}

	r.logger.Debug().Str("context", logContext).Str("code", logCode).Msg("Create user successful")
	if err = tx.Commit(ctx); err != nil {
		return models.DbUser{}, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return user, nil
}

// GetUserByLogin метод получения пользователя по логину
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
