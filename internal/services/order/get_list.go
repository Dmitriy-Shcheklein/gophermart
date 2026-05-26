package order

import (
	"context"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
)

func (s *Service) GetList(ctx context.Context, userID int) ([]models.RequestOrder, error) {
	dbOrders, err := s.repository.GetByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]models.RequestOrder, len(dbOrders))
	for i := range dbOrders {
		result[i] = models.RequestOrder{
			Number:     dbOrders[i].Number,
			Status:     dbOrders[i].Status,
			UploadedAt: dbOrders[i].UploadedAt,
		}
		if dbOrders[i].Accrual.Valid {
			result[i].Accrual = &dbOrders[i].Accrual.Float64
		}
	}
	return result, nil
}
