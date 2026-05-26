package bootstrap

import (
	"context"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/config"
	userHandler "github.com/Dmitriy-Shcheklein/gophermart/internal/handlers/user"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/repositories/pg"
	userSvc "github.com/Dmitriy-Shcheklein/gophermart/internal/services/user"
	"github.com/go-chi/chi/v5"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

func Bootstrap(ctx context.Context, cfg *config.Config, router *chi.Mux, logger *zerolog.Logger) error {
	pool, err := pgxpool.New(ctx, cfg.DbDSN)
	if err != nil {
		return err
	}
	repository, err := pg.New(logger, pool)
	if err != nil {
		return err
	}
	svc, err := userSvc.New(logger, repository)
	if err != nil {
		return err
	}
	handler, err := userHandler.New(logger, svc)
	if err != nil {
		return err
	}
	router.Post("/api/user/register", handler.Register)
	router.Post("/api/user/login", handler.Auth)

	return nil
}
