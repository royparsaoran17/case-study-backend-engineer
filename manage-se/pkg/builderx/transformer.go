package builderx

import (
	"time"

	"database/sql"
)

const (
	timeFormat string = "2006-01-02 15:04:05"
)

func SQLRowToMap(rows *sql.Rows) (map[string]any, error) {
	if rows == nil {
		return nil, nil
	}

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}

		return nil, nil
	}

	columnNames, err := rows.Columns()
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	values := make([]any, len(columnNames))
	valuesWrapped := make([]any, 0, len(columnNames))
	for i := range values {
		valuesWrapped = append(valuesWrapped, &values[i])
	}

	if err := rows.Scan(valuesWrapped...); err != nil {
		return nil, err
	}

	jObj := map[string]any{}
	for i, v := range values {
		col := columnNames[i]
		switch t := v.(type) {
		case string:
			jObj[col] = t
		case []byte:
			jObj[col] = string(t)
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			jObj[col] = t
		case float32, float64:
			jObj[col] = t
		case bool:
			jObj[col] = t
		case time.Time:
			jObj[col] = t.Format(timeFormat)
		default:
			jObj[col] = t
		}
	}

	return jObj, nil
}

func SQLRowsToArray(rows *sql.Rows) ([]map[string]any, error) {
	if rows == nil {
		return nil, nil
	}
	columnNames, err := rows.Columns()
	defer rows.Close()

	if err != nil {
		return nil, err
	}
	jArray := []map[string]any{}
	for rows.Next() {
		values := make([]any, len(columnNames))
		valuesWrapped := make([]any, 0, len(columnNames))
		for i := range values {
			valuesWrapped = append(valuesWrapped, &values[i])
		}

		if err := rows.Scan(valuesWrapped...); err != nil {
			return nil, err
		}

		jObj := map[string]any{}
		for i, v := range values {
			col := columnNames[i]
			switch t := v.(type) {
			case string:
				jObj[col] = t
			case []byte:
				jObj[col] = string(t)
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				jObj[col] = t
			case float32, float64:
				jObj[col] = t
			case bool:
				jObj[col] = t
			case time.Time:
				jObj[col] = t.Format(timeFormat)
			default:
				jObj[col] = t
			}
		}
		jArray = append(jArray, jObj)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jArray, nil
}
