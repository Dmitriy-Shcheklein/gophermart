package order

import (
	"context"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
)

// GetWithdrawals получение списаний баллов
func (s *Service) GetWithdrawals(ctx context.Context, userID int) ([]models.ResponseWithdrawn, error) {
	dbWithdrawals, err := s.repository.GetWithdrawals(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]models.ResponseWithdrawn, len(dbWithdrawals))
	for i := range dbWithdrawals {
		result[i] = models.ResponseWithdrawn{
			Sum:         dbWithdrawals[i].Sum,
			Order:       dbWithdrawals[i].Order,
			ProcessedAt: dbWithdrawals[i].ProcessedAt,
		}
	}
	return result, nil
}
