package builderx

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"
)

// PostgreQueryInsert gnerate postgre query insert
func PostgreQueryInsert(table string, columns []string, nVal int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "INSERT INTO %s", table)
	fmt.Fprintf(&b, "(%s) VALUES", strings.Join(columns, ", "))

	var placeholders []string
	for i := 1; i <= nVal; i++ {
		// Construct the full query by adding placeholders for each
		// set of values that we want to insert.
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		if i%len(columns) != 0 {
			continue
		}

		// When the end of a set is reached, write it to the query
		// builder and reset placeholders.
		fmt.Fprintf(&b, "(%s)", strings.Join(placeholders, ", "))
		placeholders = nil

		// Do not add a comma delimiter after the last set of values.
		if i == nVal {
			break
		}
		b.WriteString(", ")
	}

	return b.String()
}

// PostgreSqlJoin concatenates the given elements into a string.
func PostgreQueryUpdate(elems []string) string {
	b := new(strings.Builder)
	for index, e := range elems {
		b.WriteString(fmt.Sprintf("%s = $%d, ", e, index+1))
	}

	if b.Len() == 0 {
		return b.String()
	}

	return b.String()[0 : b.Len()-2]
}

func PostgreQueryWhere(elems []string, start int) string {
	b := new(strings.Builder)

	for i := 0; i < len(elems); i++ {
		b.WriteString(fmt.Sprintf("%s = $%d AND ", elems[i], i+start))
	}

	if b.Len() == 0 {
		return b.String()
	}

	return b.String()[0 : b.Len()-4]
}

func StructToPostgreQueryWhere(iStruct any, tag string) (QueryWhere, error) {

	var (
		periodRange string
		startDate,
		endDate string
		qw = QueryWhere{}
	)

	if iStruct == nil {
		return qw, nil
	}

	data, err := StructToKeyValue(iStruct, tag)
	if err != nil {
		return qw, err
	}

	if len(data) == 0 {
		return qw, err
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
		return qw, fmt.Errorf("the struct is empty value")
	}

	nw := PostgreQueryWhere(qw.Columns, 1)

	nc := len(qw.Columns)
	if len(startDate) > 0 && len(endDate) > 0 {
	    qw.Values = append(qw.Values, startDate, endDate)
		periodRange = "(created_at >= $1  AND created_at <= $2 )"
		if len(data) > 0 {
			periodRange = fmt.Sprintf("(created_at >= $%d  AND created_at <= $%d )", nc+1, nc+2)
		}
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
