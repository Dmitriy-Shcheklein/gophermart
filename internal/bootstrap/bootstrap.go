package bootstrap

import (
	"context"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/config"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/config/pgpool"
	orderHandler "github.com/Dmitriy-Shcheklein/gophermart/internal/handlers/order"
	userHandler "github.com/Dmitriy-Shcheklein/gophermart/internal/handlers/user"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/middlewares"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/repositories/pg"
	orderSvc "github.com/Dmitriy-Shcheklein/gophermart/internal/services/order"
	userSvc "github.com/Dmitriy-Shcheklein/gophermart/internal/services/user"
	"github.com/go-chi/chi/v5"

	"github.com/rs/zerolog"
)

func Bootstrap(
	_ context.Context, cfg *config.Config, router *chi.Mux, logger *zerolog.Logger, mw *middlewares.Middleware,
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
	oSvc, err := orderSvc.New(logger, repository)
	if err != nil {
		return err
	}
	authSvc := middlewares.NewAuthService()
	oHandler, err := orderHandler.New(logger, oSvc, authSvc)
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

	return nil
}
