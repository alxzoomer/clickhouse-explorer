package app

import (
	"database/sql"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/alxzoomer/clickhouse-explorer/pkg/dbexport"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Application struct {
}

func New() *Application {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if os.Getenv("APP_ENVIRONMENT") == "DEV" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	return &Application{}
}

func (app *Application) Run() {
	log.Info().Msg("Starting clickhouse-explorer. Open http://localhost:8000")
	http.HandleFunc("/", handler)
	err := http.ListenAndServe("localhost:8000", nil)
	log.Fatal().Err(err).Msg("")
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if r.RequestURI != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	j, err := queryExample()
	if err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			log.Error().Err(err).Msg("")
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(j)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func queryExample() ([]byte, error) {
	// clickhouseUrl := "tcp://127.0.0.1:9000?debug=true"
	clickhouseUrl := "tcp://127.0.0.1:9000?debug=false"
	connect, err := sql.Open("clickhouse", clickhouseUrl)
	if err != nil {
		return nil, err
	}
	defer connect.Close()
	if err := connect.Ping(); err != nil {
		return nil, err
	}

	rows, err := connect.Query("SELECT * FROM test.example_table")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return dbexport.MarshalDbRows(rows)
}
