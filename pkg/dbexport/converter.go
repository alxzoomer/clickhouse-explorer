package dbexport

import (
	"database/sql"
	"encoding/json"
)

// MarshalDbRows returns JSON representation of sql.Rows
func MarshalDbRows(rows *sql.Rows) ([]byte, error) {
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
	z, err := json.Marshal(finalRows)
	return z, err
}
