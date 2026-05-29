package order

import (
	"context"
	"sync"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
	"golang.org/x/sync/errgroup"
)

// ProcessOrders обработка не рассчитанных заказов
func (s *Service) ProcessOrders(ctx context.Context) error {
	orders, err := s.repository.GetProcessingOrders(ctx)
	s.logger.Info().Int("count", len(orders)).Msg("Processing orders")
	if err != nil {
		return err
	}
	statusesMap := map[models.LoyaltyStatuses]models.OrderStatuses{
		models.Invalid:    models.InvalidOrder,
		models.Processing: models.ProcessingOrder,
		models.Processed:  models.ProcessedOrder,
	}

	errGroup := errgroup.Group{}

	wg := sync.WaitGroup{}
	for i := range orders {
		wg.Add(1)
		go func() {
			defer wg.Done()
			orderInfo, err := s.loyaltyService.GetOrder(ctx, orders[i].Number)
			if err != nil {
				s.logger.Error().Err(err).Msg("error while get order info")
			}
			currentStatus, ok := statusesMap[orderInfo.Status]
			if !ok {
				return
			}

			orders[i].Status = currentStatus
			if orderInfo.Accrual != nil {
				orders[i].Accrual.Valid = true
				orders[i].Accrual.Float64 = *orderInfo.Accrual
			}
			if err = s.repository.UpdateOrder(ctx, orders[i]); err != nil {
				s.logger.Error().Err(err).Msg("error while update order")
			}
		}()

	}
	wg.Wait()

	return nil
}
