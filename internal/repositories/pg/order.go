package pg

import (
	"context"
	"fmt"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateOrder(ctx context.Context, userID int, orderNum string) error {
	query := "insert into orders (number, user_id, status) values($1, $2, $3)"
	_, err := r.pool.Exec(ctx, query, orderNum, userID, models.NewOrder)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetOrderByNum(ctx context.Context, orderNum string) (models.DbOrder, error) {
	query := "select id, status, uploaded_at, accrual, number, user_id from orders where number = $1"

	row := r.pool.QueryRow(ctx, query, orderNum)
	var result models.DbOrder
	if err := row.Scan(
		&result.ID, &result.Status, &result.UploadedAt, &result.Accrual, &result.Number, &result.UserID,
	); err != nil {
		return models.DbOrder{}, err
	}
	return result, nil
}

func (r *Repository) GetByUserId(ctx context.Context, userID int) ([]models.DbOrder, error) {
	query := "select id, status, uploaded_at, accrual, number from orders where user_id = $1 order by uploaded_at desc"
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	orders, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.DbOrder])
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *Repository) Withdraw(ctx context.Context, balance models.DbBalance, withdraw models.DbWithdrawn) (err error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if err = tx.Rollback(ctx); err != nil {
				r.logger.Error().Err(err).Msg("error while rollback transaction")
			}
		}
	}()
	balanceQuery := "update balances set current = $1, withdrawn = $2 where user_id = $3"
	_, err = tx.Exec(ctx, balanceQuery, balance.Current, balance.Withdrawn, balance.UserID)
	if err != nil {
		return err
	}
	withdrawQuery := "insert into withdrawns (sum, order_num, user_id) values ($1, $2, $3)"
	_, err = tx.Exec(ctx, withdrawQuery, withdraw.Sum, withdraw.Order, withdraw.UserID)
	if err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *Repository) GetUserBalance(ctx context.Context, userID int) (models.DbBalance, error) {
	r.logger.Debug().Int("userID", userID).Msg("GetUserBalance")
	query := "select id, current, withdrawn, user_id from balances where user_id = $1"

	row := r.pool.QueryRow(ctx, query, userID)
	var result models.DbBalance
	if err := row.Scan(
		&result.ID, &result.Current, &result.Withdrawn, &result.UserID,
	); err != nil {
		return models.DbBalance{}, err
	}
	r.logger.Debug().Interface("result", result).Msg("GetUserBalance - result")
	return result, nil
}

func (r *Repository) GetProcessingOrders(ctx context.Context) ([]models.DbOrder, error) {
	query := "select id, status, uploaded_at, accrual, number from orders where status != 'INVALID' and status != 'PROCESSED'"
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	orders, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.DbOrder])
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *Repository) UpdateOrder(ctx context.Context, order models.DbOrder) error {
	query := "update orders set status = $1, accrual = $2 where number = $3"

	_, err := r.pool.Exec(ctx, query, order.Status, order.Accrual, order.Number)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetWithdrawals(ctx context.Context, userID int) ([]models.DbWithdrawn, error) {
	query := "select id, sum, order_num, user_id, processed_at from withdrawns where user_id = $1 order by processed_at desc"

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	orders, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.DbWithdrawn])
	if err != nil {
		return nil, err
	}
	return orders, nil
}
