package app

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/alxzoomer/clickhouse-explorer/pkg/dbexport"
	"log"
)

type Application struct {
}

func New() *Application {
	return &Application{}
}

func (app *Application) Run() {
	fmt.Println("Starting clickhouse-explorer")
	// clickhouseUrl := "tcp://127.0.0.1:9000?debug=true"
	clickhouseUrl := "tcp://127.0.0.1:9000?debug=false"
	connect, err := sql.Open("clickhouse", clickhouseUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer connect.Close()
	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
		return
	}

	rows, err := connect.Query("SELECT * FROM test.example_table")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	js, err := dbexport.MarshalDbRows(rows)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s", string(js))
}
