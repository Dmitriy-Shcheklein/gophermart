package pg

import (
	"context"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
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
	query := "select id, status, uploaded_at, accrual, number from orders where number = $1"

	row := r.pool.QueryRow(ctx, query, orderNum)
	var result models.DbOrder
	if err := row.Scan(&result); err != nil {
		return models.DbOrder{}, err
	}
	return result, nil
}
