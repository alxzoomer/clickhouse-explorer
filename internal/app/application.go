package app

import (
	"database/sql"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/alxzoomer/clickhouse-explorer/pkg/dbexport"
	"log"
	"net/http"
)

type Application struct {
}

func New() *Application {
	return &Application{}
}

func (app *Application) Run() {
	log.Println("Starting clickhouse-explorer. Open http://localhost:8000")
	http.HandleFunc("/", handler) // each request calls handler
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
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
			log.Println(err)
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
