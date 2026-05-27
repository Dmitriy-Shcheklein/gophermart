package order

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/errors"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type Config interface {
	GetAccrualSrvAddr() string
}

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
}

type GetFunc = func(url string) (resp *http.Response, err error)

type Service struct {
	logger     *zerolog.Logger
	cfg        Config
	isWait     atomic.Bool
	mu         sync.RWMutex
	validate   *validator.Validate
	httpClient HttpClient
}

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
