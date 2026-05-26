package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Dmitriy-Shcheklein/gophermart/internal/bootstrap"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/config"
	"github.com/Dmitriy-Shcheklein/gophermart/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func main() {
	notifyCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("error while getting config: %v\n", err)
	}
	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)

	logger.Info().Str("serverAddress", cfg.GetSrvAddr()).Str("DSN", cfg.DbDsn()).Msg("config data")

	appMiddlewares, err := middlewares.New(&logger)
	if err != nil {
		log.Fatalf("error while create appMiddlewares: %v\n", err)
	}

	routerTimeout := 15 * time.Second
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(routerTimeout))
	router.Use(appMiddlewares.WithGzip)

	if err = bootstrap.RunMigration(cfg.DbDsn()); err != nil {
		log.Fatalf("error while run migrations: %v\n", err)
	}
	if err = bootstrap.Bootstrap(notifyCtx, cfg, router, &logger, appMiddlewares); err != nil {
		log.Fatalf("error while bootstrap app: %v\n", err)
	}

	server := &http.Server{
		Addr:              cfg.GetSrvAddr(),
		Handler:           router,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("http server stopped: %v", err)
		}
	}()

	log.Println("Server started")

	<-notifyCtx.Done()

	log.Printf("server shutdown started")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
	log.Printf("server shutdown successfully")
}
