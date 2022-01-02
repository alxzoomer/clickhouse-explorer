package clickhouse

import (
	"database/sql"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/alxzoomer/clickhouse-explorer/pkg/dbexport"
	"time"
)

type Model struct {
	conn *sql.DB
}

func New(dataSourceName string) (*Model, error) {
	conn, err := sql.Open("clickhouse", dataSourceName)
	if err != nil {
		return nil, err
	}
	conn.SetMaxIdleConns(20)
	conn.SetMaxOpenConns(20)
	conn.SetConnMaxIdleTime(15 * time.Minute)
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	m := &Model{
		conn: conn,
	}
	return m, nil
}

func (m *Model) Close() error {
	return m.conn.Close()
}

func (m *Model) Query(query string) ([]interface{}, error) {
	rows, err := m.conn.Query(query)
	if err != nil {
		return nil, err
	}
	return dbexport.AsArray(rows)
}
