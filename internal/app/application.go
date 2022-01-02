package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/alxzoomer/clickhouse-explorer/internal/router"
	"github.com/alxzoomer/clickhouse-explorer/pkg/db/clickhouse"
	"github.com/alxzoomer/clickhouse-explorer/pkg/logging"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

type Application struct{}

func New() *Application {
	logging.Init(os.Getenv("APP_ENVIRONMENT") == "DEV")

	return &Application{}
}

func (app *Application) Run() {
	model := app.initDbConnection()
	defer func() {
		_ = model.Close()
		log.Info().Msg("ClickHouse connection closed")
	}()

	srv := app.initHttpServer(model)

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

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}
		log.Info().Msg("graceful shutdown completed")
		shutdownError <- nil
	}()

	go func() {
		err := srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			shutdownError <- err
		}
	}()

	err := <-shutdownError
	log.Err(err).Msg("service stopped")
}

func (app *Application) initDbConnection() *clickhouse.Model {
	clickhouseUrl := "tcp://127.0.0.1:9000?debug=false"
	model, err := clickhouse.New(clickhouseUrl)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to connect to ClickHouse database")
	}
	log.Info().
		Str("url", clickhouseUrl).
		Msg("service connected to ClickHouse database")
	return model
}

func (app *Application) initHttpServer(model router.Model) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", 8000),
		Handler:      router.New(model).Handler(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     logging.NewErrorLog(),
	}
}
