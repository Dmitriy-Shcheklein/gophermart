package order

import (
	"context"
	"math"

	domainErrors "github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
)

// Withdraw Списание баллов
func (s *Service) Withdraw(ctx context.Context, userID int, sum float64, orderNum string) error {
	if !validateLuhn(orderNum) {
		return domainErrors.ErrOrderInvalidNumber
	}
	currentBalance, err := s.repository.GetUserBalance(ctx, userID)
	if err != nil {
		return err
	}

	currentBalance.Current = RoundToTwo(currentBalance.Current - sum)
	currentBalance.Withdrawn = RoundToTwo(currentBalance.Withdrawn + sum)
	if currentBalance.Current < 0 {
		return domainErrors.ErrOrderNotEnoughBalance
	}

	if err = s.repository.Withdraw(
		ctx, currentBalance, models.DbWithdrawn{Order: orderNum, Sum: sum, UserID: userID},
	); err != nil {
		return err
	}
	return nil
}

func RoundToTwo(value float64) float64 {
	return math.Round(value*100) / 100
}
