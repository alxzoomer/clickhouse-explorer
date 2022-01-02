package router

import (
	"errors"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
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
	<script type="text/javascript">
		async function execQuery() {
			const url = '/api/v1/query';
			const data = { query: 'select 1' };

			try {
			  const response = await fetch(url, {
				method: 'POST',
				body: JSON.stringify(data),
				headers: {
				  'Content-Type': 'application/json'
				}
			  });
			  const json = await response.json();
			  console.log('Success:', JSON.stringify(json));
			  document.getElementById('result').innerHTML = JSON.stringify(json);
			} catch (error) {
			  console.error('Error:', error);
			  document.getElementById('result').innerHTML = JSON.stringify(error);
			}
		}
	</script>
	<h1>Under construction</h1>
	<div>
		<button type="button" onclick="execQuery()">
			Click me to exec query.</button>
	</div>
	<div id="result" style="width:100%">
	</div>
  </body>
</html>
`

	_, err := fmt.Fprint(w, html)

	log.Err(err).
		Str("method", r.Method).
		Str("uri", r.RequestURI).
		Msg("index handler")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (rt *Router) queryHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var p struct {
		Query string `json:"query"`
	}
	err := rt.readJSON(w, r, &p)
	if err != nil {
		rt.badRequestResponse(w, r, err)
		return
	}
	rows, err := rt.model.Query(p.Query)
	if err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Info().
				Int("code", int(exception.Code)).
				Str("exception_msg", exception.Message).
				Str("exception_name", exception.Name).
				Msg("")
			// Wrap message without additional details like internal error code and name into separate error instance
			rt.badRequestResponse(w, r, errors.New(exception.Message))
			return
		}
		log.Error().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	e := envelope{"rows": rows, "count": len(rows)}
	err = rt.writeJSON(w, http.StatusOK, e, nil)
	if err != nil {
		rt.internalServerErrorResponse(w, r, err)
	}
}
