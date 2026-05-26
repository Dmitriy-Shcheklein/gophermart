package order

import (
	"context"
	"errors"

	domainErrors "github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Service) Upload(ctx context.Context, userID int, orderNum string) error {
	if !validateLuhn(orderNum) {
		return domainErrors.ErrOrderInvalidNumber
	}
	if err := s.repository.CreateOrder(ctx, userID, orderNum); err != nil {
		if isUniqueViolation(err) {
			order, gErr := s.repository.GetOrderByNum(ctx, orderNum)
			if gErr != nil {
				return gErr
			}
			if order.UserID != userID {
				return domainErrors.ErrOrderBelongsAnotherUser
			}
			return domainErrors.ErrOrderAlreadyExists
		}
		return err
	}
	return nil
}

func validateLuhn(number string) bool {
	var sum int
	var alternate bool

	for i := len(number) - 1; i >= 0; i-- {
		n := int(number[i] - '0')
		if n < 0 || n > 9 {
			return false
		}
		if alternate {
			n *= 2
			if n > 9 {
				n = n/10 + n%10
			}
		}

		sum += n
		alternate = !alternate
	}

	return sum%10 == 0
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
