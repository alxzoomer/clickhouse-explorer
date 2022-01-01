package app

import (
	"fmt"
	"github.com/alxzoomer/clickhouse-explorer/internal/router"
	"github.com/alxzoomer/clickhouse-explorer/pkg/logging"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Application struct {
	srv *http.Server
}

func New() *Application {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if os.Getenv("APP_ENVIRONMENT") == "DEV" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

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
	log.Info().Msg("Starting clickhouse-explorer. Open http://localhost:8000")
	err := app.srv.ListenAndServe()
	log.Fatal().Err(err).Msg("")
}
