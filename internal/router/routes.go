package router

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/alxzoomer/clickhouse-explorer/pkg/dbexport"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Router struct {
	routes http.Handler
}

func New() *Router {
	router := httprouter.New()
	rt := &Router{
		routes: router,
	}
	router.NotFound = http.HandlerFunc(rt.notFoundHandler)
	router.PanicHandler = rt.panicHandler
	router.GET("/", rt.indexHandler)
	router.GET("/api/v1/query", rt.queryHandler)
	return rt
}

func (rt *Router) Handler() http.Handler {
	return rt.routes
}

func (rt *Router) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Uri     string `json:"uri"`
	}{
		Status:  http.StatusNotFound,
		Message: "Not Found",
		Uri:     r.RequestURI,
	}
	js, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	if _, err = w.Write(js); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (rt *Router) panicHandler(w http.ResponseWriter, r *http.Request, rcv interface{}) {
	log.Error().
		Interface("recovery", rcv).
		Str("uri", r.RequestURI).
		Str("method", r.Method).
		Msg("Internal server error")

	data := struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Uri     string `json:"uri"`
	}{
		Status:  http.StatusInternalServerError,
		Message: "500 Internal server error",
		Uri:     r.RequestURI,
	}
	js, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	if _, err = w.Write(js); err != nil {
		log.Error().
			Err(err).
			Str("uri", r.RequestURI).
			Str("method", r.Method).
			Msg("Internal server error")
	}
}

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
