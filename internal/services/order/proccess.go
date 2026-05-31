package order

import (
	"context"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
	"golang.org/x/sync/errgroup"
)

// ProcessOrders обработка не рассчитанных заказов
func (s *Service) ProcessOrders(ctx context.Context) error {
	limit := 10
	errGroup, gCtx := errgroup.WithContext(ctx)
	errGroup.SetLimit(limit)

	for order := range s.repository.GetProcessingOrders(ctx) {
		errGroup.Go(
			func() error {
				select {
				case <-gCtx.Done():
					return gCtx.Err()
				default:
					return s.handleOrder(gCtx, order)
				}
			},
		)
	}
	if err := errGroup.Wait(); err != nil {
		s.logger.Error().Err(err).Msg("worker finished with err")
		return err
	}
	return nil
}

func (s *Service) handleOrder(ctx context.Context, order models.DbOrder) error {
	statusesMap := map[models.LoyaltyStatuses]models.OrderStatuses{
		models.Invalid:    models.InvalidOrder,
		models.Processing: models.ProcessingOrder,
		models.Processed:  models.ProcessedOrder,
	}
	orderInfo, err := s.loyaltyService.GetOrder(ctx, order.Number)
	if err != nil {
		s.logger.Error().Err(err).Msg("error while get order info")
	}
	currentStatus, ok := statusesMap[orderInfo.Status]
	if !ok {
		return nil
	}
	order.Status = currentStatus
	if orderInfo.Accrual != nil {
		order.Accrual.Valid = true
		order.Accrual.Float64 = *orderInfo.Accrual
	}
	if err = s.repository.UpdateOrder(ctx, order); err != nil {
		return err
	}
	return nil
}
