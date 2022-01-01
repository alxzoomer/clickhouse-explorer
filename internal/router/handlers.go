package router

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/alxzoomer/clickhouse-explorer/pkg/dbexport"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (rt *Router) indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	html := `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>ClickHouse Explorer</title>
    <link rel="stylesheet" href="style.css">
  </head>
  <body>
	<h1>Under construction</h1>
	<div>
		<a href="http://localhost:8000/api/v1/query">Example query</a>
	</div>
  </body>
</html>
`

	_, err := fmt.Fprintf(w, html)

	log.Err(err).
		Str("method", r.Method).
		Str("uri", r.RequestURI).
		Msg("index handler")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (rt *Router) queryHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	log.Err(err).
		Str("method", r.Method).
		Str("uri", r.RequestURI).
		Msg("index handler")
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
