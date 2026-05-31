package order

import (
	"context"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
)

// GetBalance получение баланса пользователя
func (s *Service) GetBalance(ctx context.Context, userID int) (models.ResponseBalance, error) {
	var balance models.ResponseBalance
	dbBalance, err := s.repository.GetUserBalance(ctx, userID)
	if err != nil {
		return balance, err
	}
	balance.Withdrawn = dbBalance.Withdrawn
	balance.Current = dbBalance.Current
	return balance, nil
}
