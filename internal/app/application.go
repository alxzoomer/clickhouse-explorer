package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/alxzoomer/clickhouse-explorer/internal/router"
	"github.com/alxzoomer/clickhouse-explorer/pkg/logging"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

type Application struct {
	srv *http.Server
}

func New() *Application {
	logging.Init(os.Getenv("APP_ENVIRONMENT") == "DEV")

	return &Application{
		srv: &http.Server{
			Addr:         fmt.Sprintf(":%d", 8000),
			Handler:      router.New().Handler(),
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			ErrorLog:     logging.NewErrorLog(),
		},
	}
}

func (app *Application) Run() {
	log.Info().Msg("starting clickhouse-explorer. open http://localhost:8000")
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		log.Info().
			Interface("signal", s).
			Msg("caught signal, shutting down HTTP server gracefully")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := app.srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}
		log.Info().Msg("graceful shutdown completed")
		shutdownError <- nil
	}()

	go func() {
		err := app.srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			shutdownError <- err
		}
	}()

	err := <-shutdownError
	log.Err(err).Msg("service stopped")
}
