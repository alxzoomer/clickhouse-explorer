package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ClickHouse/clickhouse-go"
)

func main() {
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

	for rows.Next() {
		var (
			id          int
			name        string
			description string
		)
		if err := rows.Scan(&id, &name, &description); err != nil {
			log.Fatal(err)
		}
		log.Printf("ID: %d; Name: %s; Description: %s", id, name, description)
	}
}
