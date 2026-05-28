package bootstrap

import (
	"context"
	"net/http"
	"time"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/config"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/config/pgpool"
	orderHandler "github.com/Dmitriy-Shcheklein/gophermart/internal/handlers/order"
	userHandler "github.com/Dmitriy-Shcheklein/gophermart/internal/handlers/user"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/middlewares"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/repositories/pg"
	loyaltySvc "github.com/Dmitriy-Shcheklein/gophermart/internal/services/loyalty"
	orderSvc "github.com/Dmitriy-Shcheklein/gophermart/internal/services/order"
	userSvc "github.com/Dmitriy-Shcheklein/gophermart/internal/services/user"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/workers/handle_orders"
	"github.com/go-chi/chi/v5"

	"github.com/rs/zerolog"
)

func Bootstrap(
	ctx context.Context, cfg *config.Config, router *chi.Mux, logger *zerolog.Logger, mw *middlewares.Middleware,
) error {
	pool, err := pgpool.NewPool(cfg.DbDsn())
	if err != nil {
		return err
	}
	repository, err := pg.New(logger, pool.Pool)
	if err != nil {
		return err
	}
	uSvc, err := userSvc.New(logger, repository)
	if err != nil {
		return err
	}
	uHandler, err := userHandler.New(logger, uSvc)
	if err != nil {
		return err
	}
	lSvc, err := loyaltySvc.New(logger, cfg, &http.Client{Timeout: 10 * time.Second})
	if err != nil {
		return err
	}
	oSvc, err := orderSvc.New(logger, repository, lSvc)
	if err != nil {
		return err
	}
	authSvc := middlewares.NewAuthService()
	oHandler, err := orderHandler.New(logger, oSvc, authSvc)
	if err != nil {
		return err
	}
	worker, err := handle_orders.New(logger, oSvc)
	if err != nil {
		return err
	}

	router.Post("/api/user/register", uHandler.Register)
	router.Post("/api/user/login", uHandler.Auth)
	router.Route(
		"/api/user/orders", func(r chi.Router) {
			r.Use(mw.Auth)
			r.Post("/", oHandler.Upload)
			r.Get("/", oHandler.GetList)
		},
	)
	router.Route(
		"/api/user/balance", func(r chi.Router) {
			r.Use(mw.Auth)
			r.Post("/withdraw", oHandler.Withdraw)
			r.Get("/", oHandler.GetBalance)
		},
	)
	router.Route(
		"/api/user/withdrawals", func(r chi.Router) {
			r.Use(mw.Auth)
			r.Get("/", oHandler.GetWithdrawals)
		},
	)

	worker.Start(ctx, 1*time.Second)

	go func() {
		<-ctx.Done()
		worker.Stop()
	}()

	return nil
}
