package pg

import (
	"context"

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
