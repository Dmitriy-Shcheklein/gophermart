// Package worker пакет, который реализует воркер для фоновой обработки статусов заказов
package worker

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/rs/zerolog"
)

// Service интерфейс службы для реализации логики обработки
type Service interface {
	ProcessOrders(ctx context.Context) error
}

// Worker структура воркера
type Worker struct {
	logger     *zerolog.Logger
	service    Service
	isStarted  atomic.Bool
	mu         sync.Mutex
	cancelFunc context.CancelFunc
	timeout    time.Duration
}

// New конструктор воркера
func New(logger *zerolog.Logger, service Service) (*Worker, error) {
	if logger == nil {
		return nil, fmt.Errorf("%w: logger", errors.ErrEmptyDep)
	}
	if service == nil {
		return nil, fmt.Errorf("%w: service", errors.ErrEmptyDep)
	}
	return &Worker{logger: logger, service: service}, nil
}

// Start запуск воркера
func (w *Worker) Start(ctx context.Context, timeout time.Duration) {
	if w.isStarted.Load() {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.isStarted.Load() {
		return
	}
	witchCancel, cancelFunc := context.WithCancel(ctx)
	w.cancelFunc = cancelFunc
	if timeout > 0 {
		w.timeout = timeout
	}

	go w.handle(witchCancel)

	w.isStarted.Store(true)
	w.logger.Info().Msg("worker started")
}

// Stop остановка воркера
func (w *Worker) Stop() {
	if !w.isStarted.Load() {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.isStarted.Load() {
		return
	}

	if w.cancelFunc != nil {
		w.cancelFunc()
	}
	w.isStarted.Store(false)
}

func (w *Worker) handle(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := w.service.ProcessOrders(ctx); err != nil {
				w.logger.Error().Err(err).Msg("Error while handling order statuses")
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(w.timeout):
			}
		}
	}
}
