// Package loyalty пакет для получения данных по расчету баллов
package loyalty

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

// Config интерфейс конфигурации
type Config interface {
	// GetAccrualSrvAddr метод получения адреса системы расчета баллов
	GetAccrualSrvAddr() string
}

// HttpClient интерфейс клиента
type HttpClient interface {
	// Get метода для формирования GET запроса
	Get(url string) (resp *http.Response, err error)
}

// Service структура описывающая сервис
type Service struct {
	logger     *zerolog.Logger
	cfg        Config
	isWait     atomic.Bool
	mu         sync.RWMutex
	validate   *validator.Validate
	httpClient HttpClient
}

// New конструктор сервиса
func New(logger *zerolog.Logger, cfg Config, httpClient HttpClient) (*Service, error) {
	if logger == nil {
		return nil, fmt.Errorf("%w: logger", errors.ErrEmptyDep)
	}
	if cfg == nil {
		return nil, fmt.Errorf("%w: config", errors.ErrEmptyDep)
	}
	if httpClient == nil {
		return nil, fmt.Errorf("%w: httpClient", errors.ErrEmptyDep)
	}
	return &Service{
		logger: logger, cfg: cfg, validate: validator.New(validator.WithRequiredStructEnabled()),
		httpClient: httpClient,
	}, nil
}
