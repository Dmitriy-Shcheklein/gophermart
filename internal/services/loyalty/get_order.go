package loyalty

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	domainErrors "github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/models"
)

// GetOrder метода получения данных по заказу
func (s *Service) GetOrder(ctx context.Context, orderNum string) (models.LoyaltyOrderData, error) {
	var order models.LoyaltyOrderData
	if s.isWait.Load() {
		return order, domainErrors.ErrorLoyaltyWait
	}
	addr := s.cfg.GetAccrualSrvAddr() + "/api/orders/" + orderNum
	resp, err := s.httpClient.Get(addr)
	if err != nil {
		return order, err
	}
	defer func() {
		if resp.Body != nil {
			if err = resp.Body.Close(); err != nil {
				s.logger.Error().Err(err).Msg("error while close body")
			}
		}
	}()
	if resp.StatusCode == http.StatusTooManyRequests {
		retryHeader := resp.Header.Get("Retry-After")
		duration, cErr := strconv.Atoi(retryHeader)
		if cErr != nil {
			return order, cErr
		}

		s.waitOpen(ctx, time.Duration(duration))
		return order, domainErrors.ErrorLoyaltyTooManyRequest
	}
	if resp.StatusCode == http.StatusNoContent {
		order.Order = orderNum
		order.Status = models.Invalid
		return order, nil
	}
	if resp.StatusCode == http.StatusInternalServerError {
		return order, domainErrors.ErrLoyaltyUnknown
	}
	if resp.StatusCode != http.StatusOK {
		return order, domainErrors.ErrLoyaltyUnknownStatusCode
	}
	if err = json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return models.LoyaltyOrderData{}, fmt.Errorf("%w: %w: ", domainErrors.ErrLoyaltyDecodeBody, err)
	}
	if err = s.validate.Struct(order); err != nil {
		return models.LoyaltyOrderData{}, fmt.Errorf("%w: %w: ", domainErrors.ErrLoyaltyValidateBody, err)
	}
	return order, nil
}

func (s *Service) waitOpen(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
		return
	default:
		if !s.isWait.Load() {
			s.mu.Lock()
			defer s.mu.Unlock()
			if !s.isWait.Load() {
				s.isWait.Store(true)
			} else {
				return
			}
			go func() {
				select {
				case <-ctx.Done():
					return
				case <-time.After(time.Second * duration):
					s.isWait.Store(false)
				}
			}()

		}
	}
}
