package dbexport

import (
	"database/sql"
)

// AsArray returns structured representation of sql.Rows
func AsArray(rows *sql.Rows) ([]interface{}, error) {
	columns, err := rows.Columns()

	if err != nil {
		return nil, err
	}

	count := len(columns)
	var finalRows []interface{}

	for rows.Next() {

		scanArgs := make([]interface{}, count)
		for i := range columns {
			scanArgs[i] = new(interface{})
		}

		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		masterData := map[string]interface{}{}
		for i, v := range columns {
			masterData[v] = scanArgs[i]
		}

		finalRows = append(finalRows, masterData)
	}
	if finalRows == nil {
		// In case when no data in rows the function returns empty array
		finalRows = []interface{}{}
	}
	return finalRows, err
}
