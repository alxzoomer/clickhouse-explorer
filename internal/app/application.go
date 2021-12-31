package app

import (
	"github.com/alxzoomer/clickhouse-explorer/internal/router"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Application struct {
	router *router.Router
}

func New() *Application {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if os.Getenv("APP_ENVIRONMENT") == "DEV" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return &Application{
		router: router.New(),
	}
}

func (app *Application) Run() {
	log.Info().Msg("Starting clickhouse-explorer. Open http://localhost:8000")
	err := http.ListenAndServe("localhost:8000", app.router.Handler())
	log.Fatal().Err(err).Msg("")
}
