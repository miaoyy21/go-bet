package xmd

import (
	"database/sql"
	"strconv"
)

func sqlQuery(db *sql.DB, query string, args ...interface{}) ([]map[string]int, error) {

	// Rows
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Columns
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]sql.RawBytes, len(columns))
	valuePts := make([]interface{}, len(columns))
	for i := 0; i < len(columns); i++ {
		valuePts[i] = &values[i]
	}

	entries := make([]map[string]int, 0)
	for rows.Next() {
		if err := rows.Scan(valuePts...); err != nil {
			return nil, err
		}

		entry := make(map[string]int)
		for i, col := range columns {
			if values[i] != nil {
				sValue := string(values[i])
				iValue, err := strconv.Atoi(sValue)
				if err != nil {
					return nil, err
				}

				entry[col] = iValue
			}
		}

		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
