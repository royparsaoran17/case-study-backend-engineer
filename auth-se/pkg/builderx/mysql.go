package builderx

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/cast"

	"auth-se/pkg/util"
)

// StructToMySqlQueryWhere create query builder from struct
func StructToMySqlQueryWhere(iStruct any, tag string) (qw QueryWhere, err error) {

	var (
		startDate,
		endDate,
		periodRange string
	)

	if iStruct == nil {
		return qw, fmt.Errorf("input struct is %v", iStruct)
	}

	data, err := StructToKeyValue(iStruct, tag)
	if err != nil {
		return qw, err
	}

	if len(data) == 0 {
		return qw, fmt.Errorf("input data is %d", len(data))
	}

	for i := 0; i < len(data); i++ {
		if data[i].Key == "page" {
			qw.Page = cast.ToInt64(data[i].Value)
			continue
		}

		if data[i].Key == "limit" {
			qw.Limit = cast.ToInt64(data[i].Value)
			continue
		}

		if data[i].Key == "start_date" {
			startDate = cast.ToString(data[i].Value)
			continue
		}

		if data[i].Key == "end_date" {
			endDate = cast.ToString(data[i].Value)
			continue
		}

		qw.Values = append(qw.Values, data[i].Value)
		qw.Columns = append(qw.Columns, data[i].Key)
	}

	if (len(startDate) > 0 && len(endDate) == 0) || (len(startDate) == 0 && len(endDate) > 0) {
		return qw, fmt.Errorf("invalid date period start %s end %s", startDate, endDate)
	}

	if len(data) == 0 && (len(startDate) == 0 && len(endDate) == 0) {
		return qw, nil
	}

	nw := util.StringJoin(qw.Columns, "=? AND ", "=?")

	if len(startDate) > 0 && len(endDate) > 0 {
		qw.Values = append(qw.Values, startDate, endDate)
		periodRange = "(created_at >= ?  AND created_at <= ? )"
	}

	qw.Query = fmt.Sprintf(`WHERE %s %s`, nw, periodRange)

	if len(periodRange) > 0 && len(nw) > 0 {
		qw.Query = fmt.Sprintf(`WHERE %s AND %s`, nw, periodRange)
	}

	if len(qw.Values) == 0 && len(periodRange) == 0 {
		qw.Query = ""
	}

	return qw, err
}

// StructToQueryInsert struct to query insert builder
// this method doesn't support nested struct
func StructToQueryInsert(s any, tableName, tag string) (string, []any, error) {

	cols, vals, err := ToColumnsValues(s, tag)
	if err != nil {
		return "", vals, err
	}

	q := `INSERT INTO %s (%s) VALUES(%s)`

	if len(cols) == 0 {
		return "", nil, fmt.Errorf("no column available")
	}

	pattern := "?" + strings.Repeat(", ?", len(cols)-1)

	q = fmt.Sprintf(q, tableName, strings.Join(cols, ","), pattern)

	return q, vals, nil
}

// StructToQueryDelete struct to query delete builder
// this method doesn't support nested struct
func StructToQueryDelete(where any, tableName, tag string, soft bool) (string, []any, error) {

	cols, vals, err := ToColumnsValues(where, tag)
	if err != nil {
		return "", vals, err
	}

	q := fmt.Sprintf(`DELETE FROM %s `, tableName)

	if soft {
		q = fmt.Sprintf(`UPDATE %s SET deleted_at = ?`, tableName)
		vals = append([]any{time.Now().Format(layoutDateTimeFormat)}, vals...)
	}

	if len(cols) > 0 {
		q = fmt.Sprintf(`%s WHERE %s`, q, util.StringJoin(cols, "=?, ", "=?"))
	}

	return q, vals, nil
}

func SliceStructToBulkInsert(src any, tag string) ([]string, []any, []string, error) {
	var columns []string
	var replacers []string
	var values []any

	v := reflect.Indirect(reflect.ValueOf(src))
	t := reflect.TypeOf(src)
	if t.Kind() == reflect.Ptr {
		//v = v.Elem()
		t = t.Elem()
	}

	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return columns, values, replacers, fmt.Errorf("only accepted %s, got %s", reflect.Struct.String(), t.Kind().String())
	}

	for i := 0; i < v.Len(); i++ {

		item := v.Index(i)
		if !item.IsValid() {
			continue
		}

		cols, val, err := ToColumnsValues(item.Interface(), tag)
		if err != nil {
			return columns, values, replacers, err
		}

		if len(columns) == 0 {
			columns = cols
		}

		pattern := fmt.Sprintf(`(%s)`, strings.TrimRight(strings.Repeat("?,", len(columns)), `,`))
		replacers = append(replacers, pattern)
		values = append(values, val...)
	}

	return columns, values, replacers, nil
}

// StructToQueryUpdate struct to query update builder
// this method doesn't support nested struct
func StructToQueryUpdate(input any, where any, tableName, tag string) (string, []any, error) {

	cols, vals, err := ToColumnsValues(input, tag)
	if err != nil {
		return "", vals, err
	}

	cu, vu, err := ToColumnsValues(where, tag)
	if err != nil {
		return "", vals, err
	}

	q := fmt.Sprintf(`UPDATE %s SET %s`, tableName, util.StringJoin(cols, "=?, ", "=?"))
	if len(cu) > 0 {
		q = fmt.Sprintf(`%s WHERE %s`, q, util.StringJoin(cu, "=? AND ", "=?"))
		vals = append(vals, vu...)
	}

	return q, vals, nil
}
